// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
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

	acquire := func(acquireCtx context.Context) (int, error) {
		if acquireCtx != ctx {
			t.Fatal("acquireCtx != ctx", acquireCtx, ctx)
		}

		atomic.AddInt64(&acquireLimit, 1)
		atomic.AddInt64(&releaseLimit, 1)
		return 0, nil
	}

	release := func(releaseCtx context.Context, resource int) error {
		if releaseCtx != ctx {
			t.Fatal("releaseCtx != ctx", releaseCtx, ctx)
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

			if status.Active > pool.limit {
				t.Errorf("status.Active %d > pool.limit %d", status.Active, pool.limit)
				return
			}

			if status.Idle > pool.limit {
				t.Errorf("status.Idle %d > pool.limit %d", status.Idle, pool.limit)
				return
			}

			time.Sleep(time.Second)
		}
	}()

	var totalWaited1 = 1024
	for i := 0; i < totalWaited1; i++ {
		resource, err := pool.Take(ctx)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(5 * time.Millisecond)
		pool.Put(ctx, resource)

		status := pool.Status()
		if status.Active != 1 {
			t.Fatalf("status.Active %d is wrong", status.Active)
		}

		if status.Idle != 1 {
			t.Fatalf("status.Idle %d is wrong", status.Idle)
		}

		if status.AverageWaitDuration <= 0 {
			t.Fatal("status.AverageWaitDuration is wrong")
		}
	}

	t.Logf("%+v", pool.Status())

	if pool.totalWaited <= 0 {
		t.Fatal("pool.totalWaited is wrong")
	}

	if pool.totalWaitedDuration <= 0 {
		t.Fatal("pool.totalWaitedDuration is wrong")
	}

	var totalWaited2 = 65536
	var wg sync.WaitGroup
	for i := 0; i < totalWaited2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			resource, err := pool.Take(ctx)
			if err != nil {
				t.Error(err)
				return
			}

			time.Sleep(5 * time.Millisecond)
			pool.Put(ctx, resource)

			status := pool.Status()
			if status.Active > pool.limit {
				t.Errorf("status.Active %d > pool.limit %d", status.Active, pool.limit)
				return
			}

			if status.Idle > pool.limit {
				t.Errorf("status.Idle %d > pool.limit %d", status.Idle, pool.limit)
				return
			}

			if status.Waiting > 0 && pool.totalWaited > 0 && status.AverageWaitDuration <= 0 {
				t.Errorf("status.AverageWaitDuration %d is wrong", status.AverageWaitDuration)
				return
			}
		}()
	}

	wg.Wait()
	t.Logf("%+v", pool.Status())

	totalWaited := uint64(totalWaited1 + totalWaited2)
	if pool.totalWaited > totalWaited {
		t.Fatalf("pool.totalWaited %d > totalWaited %d", pool.totalWaited, totalWaited)
	}

	if pool.totalWaited > 0 && pool.totalWaitedDuration <= 0 {
		t.Fatalf("pool.totalWaitedDuration %d is wrong", pool.totalWaitedDuration)
	}
}
