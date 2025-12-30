// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	errPoolClosed = errors.New("rego: pool is closed")
)

// AcquireFunc is a function acquires a new resource and returns error if failed.
type AcquireFunc[Resource any] func(ctx context.Context) (Resource, error)

// ReleaseFunc is a function releases a resource and returns error if failed.
type ReleaseFunc[Resource any] func(ctx context.Context, resource Resource) error

// AvailableFunc is a function checks if a resource is available.
type AvailableFunc[Resource any] func(ctx context.Context, resource Resource) bool

// PoolClosedErrFunc is a function returns a pool closed error.
type PoolClosedErrFunc func(ctx context.Context) error

// Pool stores some resources and you can reuse them.
type Pool[Resource any] struct {
	resources chan Resource
	closed    bool

	acquire      AcquireFunc[Resource]
	release      ReleaseFunc[Resource]
	available    AvailableFunc[Resource]
	newClosedErr PoolClosedErrFunc

	limit          uint64
	active         uint64
	waiting        uint64
	waited         uint64
	waitedDuration time.Duration

	lock sync.RWMutex
}

// New returns a new pool with limit resources.
// All resources are acquired by acquire function and released by release function.
func New[Resource any](limit uint64, acquire AcquireFunc[Resource], release ReleaseFunc[Resource]) *Pool[Resource] {
	if limit <= 0 {
		panic("rego: limit <= 0")
	}

	if acquire == nil || release == nil {
		panic("rego: acquire function or release function is nil")
	}

	available := func(context.Context, Resource) bool { return true }
	newClosedErr := func(context.Context) error { return errPoolClosed }

	pool := &Pool[Resource]{
		limit:        limit,
		acquire:      acquire,
		release:      release,
		available:    available,
		newClosedErr: newClosedErr,
		resources:    make(chan Resource, limit),
		closed:       false,
	}

	return pool
}

func (p *Pool[Resource]) WithAvailableFunc(available AvailableFunc[Resource]) *Pool[Resource] {
	if available != nil {
		p.lock.Lock()
		p.available = available
		p.lock.Unlock()
	}

	return p
}

func (p *Pool[Resource]) WithPoolClosedErrFunc(newClosedErr PoolClosedErrFunc) *Pool[Resource] {
	if newClosedErr != nil {
		p.lock.Lock()
		p.newClosedErr = newClosedErr
		p.lock.Unlock()
	}

	return p
}

func (p *Pool[Resource]) acquireIdle() (resource Resource, ok bool) {
	select {
	case resource := <-p.resources:
		return resource, true
	default:
		return resource, false
	}
}

func (p *Pool[Resource]) waitIdle(ctx context.Context) (resource Resource, err error) {
	select {
	case resource := <-p.resources:
		return resource, nil
	case <-ctx.Done():
		return resource, ctx.Err()
	}
}

// Acquire acquires a resource from pool and returns an error if failed.
// You should call Pool.Release to return the resource back to the pool.
func (p *Pool[Resource]) Acquire(ctx context.Context) (resource Resource, err error) {
	for {
		p.lock.Lock()
		if p.closed {
			p.lock.Unlock()

			err = p.newClosedErr(ctx)
			return resource, err
		}

		// Try to acquire a idle resource from pool.
		var ok bool
		if resource, ok = p.acquireIdle(); ok {
			p.lock.Unlock()

			if p.available(ctx, resource) {
				return resource, nil
			}

			p.lock.Lock()
			p.active--
			p.lock.Unlock()

			if err = p.release(ctx, resource); err != nil {
				return resource, err
			}

			continue
		}

		// No idle resource, we should acquire a new one or wait a idle one.
		// Increase the active and unlock here may cause the pool becomes exhausted in advance.
		// However, we think this is acceptable in most situations.
		if p.active < p.limit {
			p.active++
			p.lock.Unlock()

			resource, err = p.acquire(ctx)
			if err != nil {
				p.lock.Lock()
				p.active--
				p.lock.Unlock()
			}

			return resource, err
		}

		p.waiting++
		p.lock.Unlock()

		startTime := time.Now()
		resource, err = p.waitIdle(ctx)
		endTime := time.Now()

		p.lock.Lock()
		p.waiting--
		p.waited++
		p.waitedDuration += endTime.Sub(startTime)
		p.lock.Unlock()

		if err != nil {
			return resource, err
		}

		if p.available(ctx, resource) {
			return resource, nil
		}

		p.lock.Lock()
		p.active--
		p.lock.Unlock()

		if err = p.release(ctx, resource); err != nil {
			return resource, err
		}
	}
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

			p.active--
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
