// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
	"sync"
	"time"
)

// AcquireFunc is a function acquires a new resource and returns error if failed.
type AcquireFunc[T any] func(ctx context.Context) (T, error)

// ReleaseFunc is a function releases a resource and returns error if failed.
type ReleaseFunc[T any] func(ctx context.Context, value T) error

// Pool stores some resources and you can reuse them.
type Pool[T any] struct {
	conf config

	acquire AcquireFunc[T]
	release ReleaseFunc[T]

	resourcePool *sync.Pool
	resources    chan *resource[T]
	done         chan struct{}
	closed       bool

	limit          uint64
	active         uint64
	waiting        uint64
	waited         uint64
	waitedDuration time.Duration

	lock sync.RWMutex
}

func New[T any](limit uint64, acquire AcquireFunc[T], release ReleaseFunc[T], opts ...Option) *Pool[T] {
	if limit <= 0 {
		panic("rego: limit <= 0")
	}

	if acquire == nil || release == nil {
		panic("rego: acquire or release function is nil")
	}

	conf := newConfig().apply(opts...)

	resourcePool := &sync.Pool{
		New: func() any {
			return new(resource[T])
		},
	}

	pool := &Pool[T]{
		conf:         *conf,
		limit:        limit,
		acquire:      acquire,
		release:      release,
		resourcePool: resourcePool,
		resources:    make(chan *resource[T], limit),
		done:         make(chan struct{}),
		closed:       false,
	}

	return pool
}

func (p *Pool[T]) freeResource(resource *resource[T]) {
	resource.reset()
	p.resourcePool.Put(resource)
}

func (p *Pool[T]) newResource(value T) *resource[T] {
	resource := p.resourcePool.Get().(*resource[T])
	resource.value = value
	return resource
}

func (p *Pool[T]) Put(ctx context.Context, value T) error {
	p.lock.RLock()
	closed := p.closed
	p.lock.RUnlock()

	if closed {
		return p.release(ctx, value)
	}

	resource := p.newResource(value)

	select {
	case p.resources <- resource:
		return nil
	case <-p.done:
		return p.release(ctx, value)
	default:
		return p.release(ctx, value)
	}
}

func (p *Pool[T]) tryToTake() (value T, ok bool) {
	select {
	case resource := <-p.resources:
		return resource.value, true
	case <-p.done:
		return value, false
	default:
		return value, false
	}
}

// waitToTake waits to take an idle resource from pool.
// Record: Add ctx.Done() to select will cause a performance problem...
// The select will call runtime.selectgo if there are more than one case in it, and runtime.selectgo has two steps which is slow:
//
//	sellock(scases, lockorder)
//	sg := acquireSudog()
//
// We don't know how to solve it yet, but we think timeout mechanism should be supported even we haven't solved it.
func (p *Pool[T]) waitToTake(ctx context.Context) (value T, err error) {
	select {
	case resource := <-p.resources:
		return resource.value, nil
	case <-ctx.Done():
		return value, ctx.Err()
	case <-p.done:
		err = p.conf.newPoolClosedErr(ctx)
		return value, err
	}
}

// Take takes a resource from pool and returns an error if failed.
// You should call pool.Put to put a taken resource back to the pool.
// We recommend you to use a defer for putting a resource safely.
func (p *Pool[T]) Take(ctx context.Context) (value T, err error) {
	p.lock.RLock()
	closed := p.closed
	p.lock.RUnlock()

	if closed {
		err = p.conf.newPoolClosedErr(ctx)
		return value, err
	}

	var ok bool
	if value, ok = p.tryToTake(); ok {
		return value, nil
	}

	p.lock.Lock()

	if p.closed {
		err = p.conf.newPoolClosedErr(ctx)
		return value, err
	}

	if p.active < p.limit {
		// Increase the active and unlock here may cause the pool becomes exhausted in advance.
		// However, we think this is acceptable in most situations.
		p.active++
		p.lock.Unlock()

		defer func() {
			if err != nil {
				p.lock.Lock()
				p.active--
				p.lock.Unlock()
			}
		}()

		return p.acquire(ctx)
	}

	startTime := time.Now()

	p.waiting++
	p.lock.Unlock()

	defer func() {
		d := time.Since(startTime)

		p.lock.Lock()
		p.waiting--
		p.waited++
		p.waitedDuration += d
		p.lock.Unlock()
	}()

	return p.waitToTake(ctx)
}

// Status returns the statistics of the pool.
func (p *Pool[T]) Status() Status {
	p.lock.RLock()
	defer p.lock.RUnlock()

	var waitDuration time.Duration
	if p.waited > 0 {
		waitDuration = p.waitedDuration / time.Duration(p.waited)
	}

	idle := uint64(len(p.resources))

	status := Status{
		Limit:        p.limit,
		Using:        p.active - idle,
		Idle:         idle,
		Waiting:      p.waiting,
		WaitDuration: waitDuration,
	}

	return status
}

func (p *Pool[T]) releaseResources(ctx context.Context) error {
	for {
		select {
		case resource := <-p.resources:
			if err := p.release(ctx, resource.value); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

// Close closes pool and releases all resources.
func (p *Pool[T]) Close(ctx context.Context) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.closed {
		return nil
	}

	if err := p.releaseResources(ctx); err != nil {
		return err
	}

	p.active = 0
	p.waiting = 0
	p.waited = 0
	p.waitedDuration = 0
	p.closed = true

	close(p.done)
	return nil
}
