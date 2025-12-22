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

// New returns a new pool with limit resources.
// All new resources are acquired by acquire function and released by release function.
func New[T any](limit uint64, acquire AcquireFunc[T], release ReleaseFunc[T], opts ...Option) *Pool[T] {
	if limit <= 0 {
		panic("rego: limit <= 0")
	}

	if acquire == nil || release == nil {
		panic("rego: acquire function or release function is nil")
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

func (p *Pool[T]) acquireIdle() (value T, ok bool) {
	select {
	case resource := <-p.resources:
		return resource.value, true
	case <-p.done:
		return value, false
	default:
		return value, false
	}
}

func (p *Pool[T]) waitIdle(ctx context.Context) (value T, err error) {
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

// Acquire acquires a resource from pool and returns an error if failed.
// You should call Pool.Release to return the resource back to the pool.
func (p *Pool[T]) Acquire(ctx context.Context) (value T, err error) {
	p.lock.Lock()
	if p.closed {
		p.lock.Unlock()

		err = p.conf.newPoolClosedErr(ctx)
		return value, err
	}

	// Try to acquire a idle resource from pool.
	var ok bool
	if value, ok = p.acquireIdle(); ok {
		p.lock.Unlock()

		return value, nil
	}

	// No idle resource, we should acquire a new one or wait a idle one.
	// Increase the active and unlock here may cause the pool becomes exhausted in advance.
	// However, we think this is acceptable in most situations.
	if p.active < p.limit {
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

	p.waiting++
	p.lock.Unlock()

	startTime := time.Now()
	defer func() {
		d := time.Since(startTime)

		p.lock.Lock()
		p.waiting--
		p.waited++
		p.waitedDuration += d
		p.lock.Unlock()
	}()

	return p.waitIdle(ctx)
}

// Release releases a resource to pool so we can reuse it next time.
func (p *Pool[T]) Release(ctx context.Context, value T) error {
	p.lock.RLock()
	if p.closed {
		p.lock.RUnlock()

		return p.release(ctx, value)
	}

	resource := p.newResource(value)

	select {
	case p.resources <- resource:
		p.lock.RUnlock()

		return nil
	case <-p.done:
		p.lock.RUnlock()

		return p.release(ctx, value)
	default:
		p.lock.RUnlock()

		return p.release(ctx, value)
	}
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

func (p *Pool[T]) releaseAll(ctx context.Context) error {
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

	if err := p.releaseAll(ctx); err != nil {
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
