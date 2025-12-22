// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// go test -v -cover -run=^TestPool$
func TestPool(t *testing.T) {
	ctx := context.Background()

	limit := int64(64)
	acquireLimit := int64(0)
	releaseLimit := int64(0)
	values := make(map[int]struct{}, 1024)

	acquire := func(acquireCtx context.Context) (int, error) {
		if acquireCtx != ctx {
			t.Fatalf("acquireCtx %p != ctx %p", acquireCtx, ctx)
		}

		atomic.AddInt64(&acquireLimit, 1)
		atomic.AddInt64(&releaseLimit, 1)

		value := rand.Int()
		values[value] = struct{}{}
		return value, nil
	}

	release := func(releaseCtx context.Context, value int) error {
		if releaseCtx != ctx {
			t.Fatalf("releaseCtx %p != ctx %p", releaseCtx, ctx)
		}

		if _, ok := values[value]; !ok {
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
