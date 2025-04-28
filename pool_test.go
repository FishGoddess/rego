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

			if status.Acquired > pool.limit {
				t.Errorf("status.Acquired %d is wrong", status.Acquired)
				return
			}

			if status.Idle > pool.limit {
				t.Errorf("status.Idle %d is wrong", status.Idle)
				return
			}

			time.Sleep(time.Second)
		}
	}()

	for i := 0; i < 1024; i++ {
		resource, err := pool.Take(ctx)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(5 * time.Millisecond)
		pool.Put(ctx, resource)

		status := pool.Status()
		if status.Acquired != 1 {
			t.Fatalf("status.Acquired %d is wrong", status.Acquired)
		}

		if status.Idle != 1 {
			t.Fatalf("status.Idle %d is wrong", status.Idle)
		}
	}

	t.Logf("%+v", pool.Status())

	var wg sync.WaitGroup
	for i := 0; i < 65536; i++ {
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
			if status.Acquired > pool.limit {
				t.Errorf("status.Acquired %d is wrong", status.Acquired)
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
}
