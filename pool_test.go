// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// go test -v -cover -run=^TestNew$
func TestNew(t *testing.T) {
	t.Run("limit_panic", func(tt *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				tt.Fatal("limit not panic")
			}
		}()

		acquire := func(context.Context) (int, error) { return 0, nil }
		release := func(context.Context, int) error { return nil }
		New(0, acquire, release)
	})

	t.Run("acquire_panic", func(tt *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				tt.Fatal("limit not panic")
			}
		}()

		release := func(context.Context, int) error { return nil }
		New(1, nil, release)
	})

	t.Run("release_panic", func(tt *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				tt.Fatal("limit not panic")
			}
		}()

		acquire := func(context.Context) (int, error) { return 0, nil }
		New(1, acquire, nil)
	})

	t.Run("not_panic", func(tt *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				tt.Fatal(r)
			}
		}()

		ctx := context.Background()

		acquire := func(context.Context) (int, error) { return 1, nil }
		release := func(context.Context, int) error { return nil }
		pool := New(1, acquire, release)
		defer pool.Close(ctx)

		value, err := pool.Acquire(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if value != 1 {
			t.Fatalf("got %+v is wrong", value)
		}

		status := pool.Status()
		if status.Using != 1 {
			t.Fatalf("got %+v is wrong", status.Using)
		}

		if status.Idle != 0 {
			t.Fatalf("got %+v is wrong", status.Idle)
		}
	})

	t.Run("acquire_idle_error", func(tt *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				tt.Fatal(r)
			}
		}()

		ctx := context.Background()

		wantErr := errors.New("wow")
		acquire := func(context.Context) (int, error) { return 0, nil }
		release := func(context.Context, int) error { return wantErr }
		available := func(context.Context, int) bool { return false }
		pool := New(1, acquire, release).WithAvailableFunc(available)
		defer pool.Close(ctx)

		err := pool.Release(ctx, 0)
		if err != nil {
			t.Fatalf("got %+v is wrong", err)
		}

		_, err = pool.Acquire(ctx)
		if err != wantErr {
			t.Fatalf("got %+v != want %+v", err, wantErr)
		}
	})

	t.Run("acquire_wait_error", func(tt *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				tt.Fatal(r)
			}
		}()

		ctx := context.Background()

		wantErr := errors.New("wow")
		acquire := func(context.Context) (int, error) { return 0, nil }
		release := func(context.Context, int) error { return wantErr }
		available := func(context.Context, int) bool { return false }
		pool := New(1, acquire, release).WithAvailableFunc(available)
		defer pool.Close(ctx)

		_, err := pool.Acquire(ctx)
		if err != nil {
			t.Fatal(err)
		}

		go func() {
			time.Sleep(time.Millisecond)
			err := pool.Release(ctx, 0)
			if err != nil {
				t.Errorf("got %+v is wrong", err)
			}
		}()

		_, err = pool.Acquire(ctx)
		if err != wantErr {
			t.Fatalf("got %+v != want %+v", err, wantErr)
		}
	})

	t.Run("acquire_wait_continue_error", func(tt *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				tt.Fatal(r)
			}
		}()

		ctx := context.Background()

		acquire := func(context.Context) (int, error) { return 0, nil }
		release := func(context.Context, int) error { return nil }
		available := func(context.Context, int) bool { return false }
		pool := New(1, acquire, release).WithAvailableFunc(available)
		defer pool.Close(ctx)

		_, err := pool.Acquire(ctx)
		if err != nil {
			t.Fatal(err)
		}

		go func() {
			time.Sleep(time.Millisecond)
			err := pool.Release(ctx, 0)
			if err != nil {
				t.Errorf("got %+v is wrong", err)
			}
		}()

		_, err = pool.Acquire(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("release_error", func(tt *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				tt.Fatal(r)
			}
		}()

		ctx := context.Background()

		wantErr := errors.New("wow")
		acquire := func(context.Context) (int, error) { return 0, wantErr }
		release := func(context.Context, int) error { return wantErr }
		pool := New(1, acquire, release)
		defer pool.Close(ctx)

		_, err := pool.Acquire(ctx)
		if err != wantErr {
			t.Fatalf("got %+v != want %+v", err, wantErr)
		}

		err = pool.Release(ctx, 0)
		if err != nil {
			t.Fatalf("got %+v is wrong", err)
		}

		err = pool.Release(ctx, 0)
		if err != wantErr {
			t.Fatalf("got %+v != want %+v", err, wantErr)
		}
	})
}

// go test -v -cover -run=^TestWithAvailableFunc$
func TestWithAvailableFunc(t *testing.T) {
	available := func(context.Context, int) bool {
		return false
	}

	pool := &Pool[int]{available: nil}
	pool.WithAvailableFunc(available)

	got := fmt.Sprintf("%p", pool.available)
	want := fmt.Sprintf("%p", available)
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}

	pool.WithAvailableFunc(nil)

	got = fmt.Sprintf("%p", pool.available)
	want = fmt.Sprintf("%p", available)
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

// go test -v -cover -run=^TestWithPoolClosedErrFunc$
func TestWithPoolClosedErrFunc(t *testing.T) {
	newClosedErr := func(context.Context) error {
		return nil
	}

	pool := &Pool[int]{newClosedErr: nil}
	pool.WithPoolClosedErrFunc(newClosedErr)

	got := fmt.Sprintf("%p", pool.newClosedErr)
	want := fmt.Sprintf("%p", newClosedErr)
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}

	pool.WithPoolClosedErrFunc(nil)

	got = fmt.Sprintf("%p", pool.newClosedErr)
	want = fmt.Sprintf("%p", newClosedErr)
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

// go test -v -cover -run=^TestPoolAcquireRelease$
func TestPoolAcquireRelease(t *testing.T) {
	ctx := context.Background()

	limit := int64(64)
	acquireLimit := int64(0)
	releaseLimit := int64(0)
	values := sync.Map{}

	acquire := func(acquireCtx context.Context) (int, error) {
		if acquireCtx != ctx {
			t.Fatalf("acquireCtx %p != ctx %p", acquireCtx, ctx)
		}

		atomic.AddInt64(&acquireLimit, 1)
		atomic.AddInt64(&releaseLimit, 1)

		value := rand.Int()
		values.Store(value, nil)
		return value, nil
	}

	release := func(releaseCtx context.Context, value int) error {
		if releaseCtx != ctx {
			t.Fatalf("releaseCtx %p != ctx %p", releaseCtx, ctx)
		}

		if _, ok := values.Load(value); !ok {
			t.Fatalf("value %d not found", value)
		}

		atomic.AddInt64(&releaseLimit, -1)
		return nil
	}

	pool := New(uint64(limit), acquire, release)
	defer func() {
		pool.Close(ctx)

		if acquireLimit != limit {
			t.Fatalf("acquireLimit %d != limit %d", acquireLimit, limit)
		}

		if releaseLimit != 0 {
			t.Fatalf("releaseLimit %d != 0", releaseLimit)
		}
	}()

	go func() {
		for {
			status := pool.Status()
			t.Logf("%+v", status)

			active := status.Using + status.Idle
			if active > pool.limit {
				t.Errorf("active %d > limit %d", active, pool.limit)
				return
			}

			if status.Idle > pool.limit {
				t.Errorf("idle %d > limit %d", status.Idle, pool.limit)
				return
			}

			time.Sleep(time.Second)
		}
	}()

	for i := 0; i < 1024; i++ {
		value, err := pool.Acquire(ctx)
		if err != nil {
			t.Fatal(err)
		}

		status := pool.Status()
		if status.Using != 1 {
			t.Fatalf("using %d is wrong", status.Using)
		}

		time.Sleep(5 * time.Millisecond)
		pool.Release(ctx, value)

		status = pool.Status()
		if status.Idle != 1 {
			t.Fatalf("idle %d is wrong", status.Idle)
		}

		if status.WaitDuration != 0 {
			t.Fatalf("status.WaitDuration %d is wrong", status.WaitDuration)
		}
	}

	t.Logf("%+v", pool.Status())

	if pool.waited != 0 {
		t.Fatalf("pool.waited %d is wrong", pool.waited)
	}

	if pool.waitedDuration != 0 {
		t.Fatalf("pool.waitedDuration %d is wrong", pool.waitedDuration)
	}

	var n = 65536
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			value, err := pool.Acquire(ctx)
			if err != nil {
				t.Error(err)
				return
			}

			status := pool.Status()
			if status.Using < 1 {
				t.Errorf("using %d is wrong", status.Using)
				return
			}

			if status.Using > pool.limit {
				t.Errorf("using %d > limit %d", status.Using, pool.limit)
				return
			}

			time.Sleep(5 * time.Millisecond)
			pool.Release(ctx, value)

			status = pool.Status()
			if status.Idle > pool.limit {
				t.Errorf("idle %d > limit %d", status.Idle, pool.limit)
				return
			}

			if status.Waiting > 0 && pool.waited > 0 && status.WaitDuration <= 0 {
				t.Errorf("wait duration %d is wrong", status.WaitDuration)
				return
			}
		}()
	}

	wg.Wait()
	t.Logf("%+v", pool.Status())

	status := pool.Status()
	if status.Using != 0 {
		t.Fatalf("using %d is wrong", status.Using)
	}

	if status.Idle != pool.limit {
		t.Fatalf("idle %d != limit %d", status.Idle, pool.limit)
	}

	if status.Waiting != 0 {
		t.Fatalf("waiting %d is wrong", status.Waiting)
	}

	if pool.waited > uint64(n) {
		t.Fatalf("pool.waited %d > n %d", pool.waited, n)
	}

	if pool.waited > 0 && pool.waitedDuration <= 0 {
		t.Fatalf("pool.waited %d > 0 but waitedDuration %d <= 0", pool.waited, pool.waitedDuration)
	}
}

// go test -v -cover -run=^TestPoolAvailable$
func TestPoolAvailable(t *testing.T) {
	ctx := context.Background()

	type Resource struct {
		available bool
	}

	acquire := func(acquireCtx context.Context) (Resource, error) {
		resource := Resource{available: true}
		return resource, nil
	}

	release := func(releaseCtx context.Context, resource Resource) error {
		return nil
	}

	available := func(ctx context.Context, resource Resource) bool {
		return resource.available
	}

	pool := New(1024, acquire, release).WithAvailableFunc(available)
	defer pool.Close(ctx)

	var n = 65536
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()

			resource, err := pool.Acquire(ctx)
			if err != nil {
				t.Error(err)
				return
			}

			if !resource.available {
				t.Errorf("resource.available %+v is wrong", resource.available)
				return
			}

			// Hava some fun :)
			if rand.IntN(n) > ii {
				resource.available = false
			}

			pool.Release(ctx, resource)
		}(i)
	}

	wg.Wait()
	t.Logf("%+v", pool.Status())
}

// go test -v -cover -run=^TestPoolClose$
func TestPoolClose(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()

	ctx := context.Background()
	releaseLimit := 0

	acquire := func(context.Context) (int, error) { return 0, nil }
	release := func(context.Context, int) error {
		releaseLimit++
		return nil
	}

	pool := New(64, acquire, release)
	defer pool.Close(ctx)

	err := pool.Close(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !pool.closed {
		t.Fatalf("got %+v is wrong", pool.closed)
	}

	_, err = pool.Acquire(ctx)
	if err != errPoolClosed {
		t.Fatalf("got %+v != want %+v", err, errPoolClosed)
	}

	err = pool.Release(ctx, 0)
	if err != nil {
		t.Fatal(err)
	}

	if releaseLimit != 1 {
		t.Fatalf("got %+v is wrong", releaseLimit)
	}
}

// go test -v -cover -run=^TestPoolTimeout$
func TestPoolTimeout(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()

	acquire := func(context.Context) (int, error) { return 0, nil }
	release := func(context.Context, int) error { return nil }

	pool := New(1, acquire, release)
	defer pool.Close(ctx)

	_, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Millisecond)
	if err = ctx.Err(); err == nil {
		t.Fatalf("got %+v is wrong", err)
	}

	_, err = pool.Acquire(ctx)
	if err != ctx.Err() {
		t.Fatalf("got %+v != want %+v", err, ctx.Err())
	}
}
