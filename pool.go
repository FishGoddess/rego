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
type AcquireFunc[Resource any] func(ctx context.Context) (Resource, error)

// ReleaseFunc is a function releases a resource and returns error if failed.
type ReleaseFunc[Resource any] func(ctx context.Context, resource Resource) error

// Pool stores some resources and you can reuse them.
type Pool[Resource any] struct {
	conf config

	acquire AcquireFunc[Resource]
	release ReleaseFunc[Resource]

	resources chan Resource
	closed    bool

	limit          uint64
	active         uint64
	waiting        uint64
	waited         uint64
	waitedDuration time.Duration

	lock sync.RWMutex
}

// New returns a new pool with limit resources.
// All resources are acquired by acquire function and released by release function.
func New[Resource any](limit uint64, acquire AcquireFunc[Resource], release ReleaseFunc[Resource], opts ...Option) *Pool[Resource] {
	if limit <= 0 {
		panic("rego: limit <= 0")
	}

	if acquire == nil || release == nil {
		panic("rego: acquire function or release function is nil")
	}

	conf := newConfig().apply(opts...)

	pool := &Pool[Resource]{
		conf:      *conf,
		limit:     limit,
		acquire:   acquire,
		release:   release,
		resources: make(chan Resource, limit),
		closed:    false,
	}

	return pool
}

func (p *Pool[Resource]) acquireIdle() (resource Resource, ok bool) {
	for {
		select {
		case resource := <-p.resources:
			return resource, true
		default:
			return resource, false
		}
	}
}

func (p *Pool[Resource]) waitIdle(ctx context.Context) (resource Resource, err error) {
	for {
		select {
		case resource := <-p.resources:
			return resource, nil
		case <-ctx.Done():
			err = ctx.Err()
			return resource, err
		}
	}
}

// Acquire acquires a resource from pool and returns an error if failed.
// You should call Pool.Release to return the resource back to the pool.
func (p *Pool[Resource]) Acquire(ctx context.Context) (resource Resource, err error) {
	p.lock.Lock()
	if p.closed {
		p.lock.Unlock()

		err = p.conf.newPoolClosedErr(ctx)
		return resource, err
	}

	// Try to acquire a idle resource from pool.
	var ok bool
	if resource, ok = p.acquireIdle(); ok {
		p.lock.Unlock()
		return resource, nil
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
func (p *Pool[Resource]) Release(ctx context.Context, resource Resource) error {
	p.lock.Lock()
	if p.closed {
		p.lock.Unlock()

		return p.release(ctx, resource)
	}

	select {
	case p.resources <- resource:
		p.lock.Unlock()

		return nil
	default:
		p.active--
		p.lock.Unlock()

		return p.release(ctx, resource)
	}
}

// Status returns the statistics of the pool.
func (p *Pool[Resource]) Status() Status {
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

func (p *Pool[Resource]) releaseAll(ctx context.Context) error {
	for {
		select {
		case resource := <-p.resources:
			if err := p.release(ctx, resource); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

// Close closes pool and releases all resources.
func (p *Pool[Resource]) Close(ctx context.Context) error {
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

	close(p.resources)
	return nil
}
