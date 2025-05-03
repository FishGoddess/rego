// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
	"testing"
	"time"
)

// go test -v -cover -run=^TestPoolStatus$
func TestPoolStatus(t *testing.T) {
	limit := 16

	pool := &Pool[int]{
		limit:               uint64(limit),
		release:             DefaultReleaseFunc[int],
		active:              4,
		waiting:             8,
		totalWaited:         0,
		totalWaitedDuration: 0,
		resources:           make(chan int, limit),
	}

	ctx := context.Background()
	defer pool.Close(ctx)

	for i := range limit {
		pool.resources <- i
	}

	poolStatus := PoolStatus{
		Limit:               pool.limit,
		Active:              pool.active,
		Idle:                uint64(len(pool.resources)),
		Waiting:             pool.waiting,
		AverageWaitDuration: 0,
	}

	if pool.Status() != poolStatus {
		t.Fatalf("pool.Status %+v != %+v", pool.Status(), poolStatus)
	}

	pool.totalWaited = 2
	pool.totalWaitedDuration = time.Second
	poolStatus.AverageWaitDuration = pool.totalWaitedDuration / time.Duration(pool.totalWaited)

	if pool.Status() != poolStatus {
		t.Fatalf("pool.Status %+v != %+v", pool.Status(), poolStatus)
	}
}
