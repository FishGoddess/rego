// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
	"sync"
	"time"

	"github.com/FishGoddess/rego/pkg/list"
	"github.com/FishGoddess/rego/pkg/token"
)

var _ ReleaseFunc[int] = DefaultReleaseFunc[int]

// DefaultReleaseFunc is a default func to release a resource.
// It does nothing to the resource.
func DefaultReleaseFunc[T any](ctx context.Context, resource T) error {
	return nil
}

// AcquireFunc is a function acquires a new resource and returns error if failed.
type AcquireFunc[T any] func(ctx context.Context) (T, error)

// ReleaseFunc is a function releases a resource and returns error if failed.
type ReleaseFunc[T any] func(ctx context.Context, resource T) error

type Pool[T any] struct {
	conf config

	acquire AcquireFunc[T]
	release ReleaseFunc[T]

	limit   uint64
	active  uint64
	waiting uint64

	totalWaited         uint64
	totalWaitedDuration time.Duration

	tokens    *token.Bucket
	resources *list.List[T]
	closed    bool

	lock sync.RWMutex
}

func New[T any](limit uint64, acquire AcquireFunc[T], release ReleaseFunc[T], opts ...Option) *Pool[T] {
	if limit <= 0 {
		panic("rego: limit can't be less than 0")
	}

	if acquire == nil || release == nil {
		panic("rego: acquire or release func can't be nil")
	}

	conf := newDefaultConfig()
	for _, opt := range opts {
		opt.ApplyTo(conf)
	}

	pool := &Pool[T]{
		conf:      *conf,
		limit:     limit,
		acquire:   acquire,
		release:   release,
		tokens:    token.NewBucket(limit),
		resources: list.New[T](),
		closed:    false,
	}

	return pool
}

func (p *Pool[T]) Put(ctx context.Context, resource T) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.closed {
		return p.release(ctx, resource)
	}

	if p.resources.Len() >= p.limit {
		return p.release(ctx, resource)
	}

	p.resources.Push(resource)
	p.produceToken()
	return nil
}

// ConsumeToken consumes a token from bucket and waits util context done if there is no token.
// Record: Add ctx.Done() to select will cause a performance problem...
// The select will call runtime.selectgo if there are more than one case in it, and runtime.selectgo has two steps which is slow:
//
//	sellock(scases, lockorder)
//	sg := acquireSudog()
//
// We don't know how to solve it yet, but we think timeout mechanism should be supported even we haven't solved it.
func (p *Pool[T]) consumeToken(ctx context.Context) (err error) {
	if p.conf.disableToken {
		return nil
	}

	p.waiting++
	p.lock.Unlock()

	startTime := time.Now()

	defer func() {
		p.lock.Lock()
		p.waiting--

		if err == nil {
			p.totalWaited++
			p.totalWaitedDuration += time.Since(startTime)
		}
	}()

	return p.tokens.ConsumeToken(ctx)
}

func (p *Pool[T]) produceToken() {
	if p.conf.disableToken {
		return
	}

	p.tokens.ProduceToken()
}

func (p *Pool[T]) tryToTake() (resource T, ok bool) {
	return p.resources.Pop()
}

// Take takes a resource from pool and returns an error if failed.
// You should call pool.Put to put a taken resource back to the pool.
// We recommend you to use a defer for putting a resource safely.
func (p *Pool[T]) Take(ctx context.Context) (resource T, err error) {
	p.lock.Lock()

	if p.closed {
		p.lock.Unlock()

		err = p.conf.newPoolClosedErr(ctx)
		return resource, err
	}

	if err := p.consumeToken(ctx); err != nil {
		p.lock.Unlock()

		return resource, err
	}

	if resource, ok := p.tryToTake(); ok {
		p.lock.Unlock()

		return resource, nil
	}

	if p.active >= p.limit {
		p.produceToken()
		p.lock.Unlock()

		err = p.conf.newPoolExhaustedErr(ctx)
		return resource, err
	}

	// Increase the active and unlock here may cause the pool becomes exhausted in advance.
	// However, we think this is acceptable in most situations.
	p.active++
	p.lock.Unlock()

	defer func() {
		if err != nil {
			p.lock.Lock()
			p.active--
			p.produceToken()
			p.lock.Unlock()
		}
	}()

	return p.acquire(ctx)
}

// Status returns the statistics of the pool.
func (p *Pool[T]) Status() PoolStatus {
	p.lock.RLock()
	defer p.lock.RUnlock()

	var averageWaitDuration time.Duration
	if p.totalWaited > 0 {
		averageWaitDuration = p.totalWaitedDuration / time.Duration(p.totalWaited)
	}

	status := PoolStatus{
		Limit:               p.limit,
		Active:              p.active,
		Idle:                p.resources.Len(),
		Waiting:             p.waiting,
		AverageWaitDuration: averageWaitDuration,
	}

	return status
}

// Close closes pool and releases all resources.
func (p *Pool[T]) Close(ctx context.Context) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.closed {
		return nil
	}

	for {
		resource, ok := p.resources.Pop()
		if !ok {
			break
		}

		if err := p.release(ctx, resource); err != nil {
			return err
		}
	}

	p.active = 0
	p.waiting = 0
	p.totalWaited = 0
	p.totalWaitedDuration = 0
	p.closed = true
	p.tokens.Free()
	return nil
}
