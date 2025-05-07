// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package token

import (
	"context"
	"testing"
	"time"
)

// go test -v -cover -run=^TestBucket$
func TestBucket(t *testing.T) {
	ctx := context.Background()
	limit := 16

	bucket := NewBucket(uint64(limit))
	defer bucket.Free()

	if len(bucket.tokens) != limit {
		t.Fatalf("len(bucket.tokens) %d != limit %d", len(bucket.tokens), limit)
	}

	for range limit {
		err := bucket.ConsumeToken(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	ctx, cancel1 := context.WithCancel(context.Background())

	go func() {
		time.Sleep(300 * time.Millisecond)
		cancel1()
	}()

	err := bucket.ConsumeToken(ctx)
	if err == nil {
		t.Fatal("bucket consume token err is nil")
	}

	if err != context.Canceled {
		t.Fatalf("bucket consume token err %v is wrong", err)
	}

	ctx, cancel2 := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel2()

	err = bucket.ConsumeToken(ctx)
	if err == nil {
		t.Fatal("bucket consume token err is nil")
	}

	if err != context.DeadlineExceeded {
		t.Fatalf("bucket consume token err %v is wrong", err)
	}

	ctx = context.Background()

	for range limit {
		bucket.ProduceToken()
	}

	if len(bucket.tokens) != limit {
		t.Fatalf("len(bucket.tokens) %d != limit %d", len(bucket.tokens), limit)
	}

	for range limit {
		err := bucket.ConsumeToken(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}
}
