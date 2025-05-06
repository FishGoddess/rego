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

	acquire := func(_ context.Context) (int, error) {
		atomic.AddInt64(&acquireLimit, 1)
		atomic.AddInt64(&releaseLimit, 1)
		return 0, nil
	}

	release := func(_ context.Context, _ int) error {
		atomic.AddInt64(&releaseLimit, -1)
		return nil
	}

	pool := New(uint64(limit), acquire, release)
	defer func() {
		if err := pool.Close(ctx); err != nil {
			t.Fatal(err)
		}

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

// go test -v -cover -run=^TestPoolExhaust$
func TestPoolContext(t *testing.T) {
	ctx := context.Background()

	limit := int64(4)

	acquire := func(acquireCtx context.Context) (int, error) {
		if acquireCtx != ctx {
			t.Fatal("acquireCtx != ctx", acquireCtx, ctx)
		}

		return 0, nil
	}

	release := func(releaseCtx context.Context, resource int) error {
		if releaseCtx != ctx {
			t.Fatal("releaseCtx != ctx", releaseCtx, ctx)
		}

		return nil
	}

	pool := New(uint64(limit), acquire, release)
	defer func() {
		if err := pool.Close(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	resources := make([]int, 0, limit)
	for range limit {
		resource, err := pool.Take(ctx)
		if err != nil {
			t.Fatal(err)
		}

		resources = append(resources, resource)
	}

	ctx, cancel1 := context.WithCancel(context.Background())

	go func() {
		time.Sleep(300 * time.Millisecond)
		cancel1()
	}()

	_, err := pool.Take(ctx)
	if err == nil {
		t.Fatal("pool take err is nil")
	}

	if err != context.Canceled {
		t.Fatalf("pool take err %v is wrong", err)
	}

	ctx, cancel2 := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel2()

	_, err = pool.Take(ctx)
	if err == nil {
		t.Fatal("pool take err is nil")
	}

	if err != context.DeadlineExceeded {
		t.Fatalf("pool take err %v is wrong", err)
	}

	for _, resource := range resources {
		if err := pool.Put(ctx, resource); err != nil {
			t.Fatal(err)
		}
	}
}

// go test -v -cover -run=^TestPoolExhaust$
func TestPoolExhaust(t *testing.T) {
	ctx := context.Background()

	limit := int64(4)
	acquireLimit := int64(0)
	releaseLimit := int64(0)

	acquire := func(_ context.Context) (int, error) {
		atomic.AddInt64(&acquireLimit, 1)
		atomic.AddInt64(&releaseLimit, 1)
		return 0, nil
	}

	release := func(_ context.Context, resource int) error {
		atomic.AddInt64(&releaseLimit, -1)
		return nil
	}

	pool := New(uint64(limit), acquire, release, WithFastFailed())
	defer func() {
		if err := pool.Close(ctx); err != nil {
			t.Fatal(err)
		}

		if acquireLimit != limit {
			t.Fatalf("acquireLimit %d != limit %d", acquireLimit, limit)
		}

		if releaseLimit != 0 {
			t.Fatalf("releaseLimit %d != 0", releaseLimit)
		}
	}()

	resources := make([]int, 0, limit)
	for range limit {
		resource, err := pool.Take(ctx)
		if err != nil {
			t.Fatal(err)
		}

		resources = append(resources, resource)
	}

	for range limit {
		_, err := pool.Take(ctx)
		if err == nil {
			t.Fatal("pool take err is nil")
		}

		if err != ErrPoolExhausted {
			t.Fatalf("pool take err %v is wrong", err)
		}
	}

	for _, resource := range resources {
		if err := pool.Put(ctx, resource); err != nil {
			t.Fatal(err)
		}
	}
}

// go test -v -cover -run=^TestPoolClose$
func TestPoolClose(t *testing.T) {
	ctx := context.Background()

	limit := int64(4)
	acquireLimit := int64(0)
	releaseLimit := int64(0)

	acquire := func(_ context.Context) (int, error) {
		atomic.AddInt64(&acquireLimit, 1)
		atomic.AddInt64(&releaseLimit, 1)
		return 0, nil
	}

	release := func(_ context.Context, _ int) error {
		atomic.AddInt64(&releaseLimit, -1)
		return nil
	}

	pool := New(uint64(limit), acquire, release)
	defer func() {
		if err := pool.Close(ctx); err != nil {
			t.Fatal(err)
		}

		if acquireLimit != limit {
			t.Fatalf("acquireLimit %d != limit %d", acquireLimit, limit)
		}

		if releaseLimit != 0 {
			t.Fatalf("releaseLimit %d != 0", releaseLimit)
		}
	}()

	resources := make([]int, 0, limit)
	for range limit {
		resource, err := pool.Take(ctx)
		if err != nil {
			t.Fatal(err)
		}

		resources = append(resources, resource)
	}

	if err := pool.Close(ctx); err != nil {
		t.Fatal(err)
	}

	_, err := pool.Take(ctx)
	if err == nil {
		t.Fatal("pool take err is nil")
	}

	if err != ErrPoolClosed {
		t.Fatalf("pool take err %v is wrong", err)
	}

	for _, resource := range resources {
		if err = pool.Put(ctx, resource); err != nil {
			t.Fatal(err)
		}
	}
}
