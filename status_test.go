// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"testing"
	"time"
)

// go test -v -cover -run=^TestStatus$
func TestStatus(t *testing.T) {
	limit := 16

	pool := &Pool[int]{
		limit:          16,
		active:         12,
		waiting:        100,
		waited:         50,
		waitedDuration: 100 * time.Millisecond,
		resources:      make(chan *resource[int], limit),
	}

	for i := range 10 {
		pool.resources <- &resource[int]{value: i}
	}

	want := Status{
		Limit:        16,
		Using:        2,
		Idle:         10,
		Waiting:      100,
		WaitDuration: 2 * time.Millisecond,
	}

	got := pool.Status()
	if got != want {
		t.Fatalf("got %+v != want %+v", got, want)
	}
}
