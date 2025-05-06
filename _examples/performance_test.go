// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"testing"

	"github.com/FishGoddess/rego"
)

// go test -v -run=^$ -bench=^BenchmarkPool$ -benchtime=1s
func BenchmarkPool(b *testing.B) {
	ctx := context.Background()

	acquire := func(ctx context.Context) (int, error) {
		return 0, nil
	}

	release := func(ctx context.Context, resource int) error {
		return nil
	}

	pool := rego.New[int](1024, acquire, release)
	defer pool.Close(ctx)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resource, err := pool.Take(ctx)
			if err != nil {
				b.Fatal(err)
			}

			if err = pool.Put(ctx, resource); err != nil {
				b.Fatal(err)
			}
		}
	})
}
