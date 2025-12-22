// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/FishGoddess/rego"
)

func acquire(ctx context.Context) (int, error) {
	return 0, nil
}

func release(ctx context.Context, resource int) error {
	return nil
}

func newPoolClosedErr(ctx context.Context) error {
	return errors.New("_examples: pool is closed")
}

func main() {
	// Create a pool with limit and new error function.
	ctx := context.Background()
	pool := rego.New(1, acquire, release, rego.WithPoolClosedErr(newPoolClosedErr))

	// Acquiring from pool is ok.
	value, err := pool.Acquire(ctx)
	fmt.Println(value, err)

	// Close the pool.
	pool.Close(ctx)

	// Now, the pool is closed so acquiring from pool will return an error.
	value, err = pool.Acquire(ctx)
	fmt.Println(value, err)
}
