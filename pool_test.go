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
				t.Errorf("status.Active %d is wrong", status.Active)
				return
			}

			if status.Idle > pool.limit {
				t.Errorf("status.Idle %d is wrong", status.Idle)
				return
			}

			time.Sleep(time.Second)
		}
	}()

	totalTaken1 := 1024
	for i := 0; i < totalTaken1; i++ {
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

		if status.AverageWaitDuration != 0 {
			t.Fatalf("status.AverageWaitDuration %d is wrong", status.AverageWaitDuration)
		}
	}

	t.Logf("%+v", pool.Status())

	if pool.totalTaken != uint64(totalTaken1) {
		t.Fatalf("pool.totalTaken %d is wrong", pool.totalTaken)
	}

	if pool.totalWaitedDuration != 0 {
		t.Fatalf("pool.totalWaitedDuration %d is wrong", pool.totalWaitedDuration)
	}

	var wg sync.WaitGroup
	totalTaken2 := 65536
	for i := 0; i < totalTaken2; i++ {
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
				t.Errorf("status.Active %d is wrong", status.Active)
				return
			}

			if status.Idle > pool.limit {
				t.Errorf("status.Idle %d is wrong", status.Idle)
				return
			}
		}()
	}

	wg.Wait()
	t.Logf("%+v", pool.Status())

	if pool.totalTaken != uint64(totalTaken1)+uint64(totalTaken2) {
		t.Fatalf("pool.totalTaken %d is wrong", pool.totalTaken)
	}

	totalTaken1 = 4096
	for i := 0; i < totalTaken1; i++ {
		resource, err := pool.Take(ctx)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Millisecond)
		pool.Put(ctx, resource)

		status := pool.Status()
		if status.Active != uint64(limit) {
			t.Fatalf("status.Active %d is wrong", status.Active)
		}

		if status.Idle != uint64(limit) {
			t.Fatalf("status.Idle %d is wrong", status.Idle)
		}

		if status.AverageWaitDuration == 0 {
			t.Fatalf("status.AverageWaitDuration %d is wrong", status.AverageWaitDuration)
		}
	}

	t.Logf("%+v", pool.Status())
}
