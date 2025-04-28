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

var (
	errPoolExhausted = errors.New("_examples: pool is exhausted")
	errPoolClosed    = errors.New("_examples: pool is closed")
)

func acquire(ctx context.Context) (int, error) {
	return 0, nil
}

func release(ctx context.Context, resource int) error {
	return nil
}

func newPoolExhaustedErr(ctx context.Context) error {
	return errPoolExhausted
}

func newPoolClosedErr(ctx context.Context) error {
	return errPoolClosed
}

func main() {
	// Create a pool with limit and fast-failed so it will return an error immediately instead of waiting.
	ctx := context.Background()
	pool := rego.New(1, acquire, release, rego.WithFastFailed())

	// Take one resource from pool which is ok.
	resource, err := pool.Take(ctx)
	fmt.Println(resource, err)

	// However, the pool is exhausted after taking one resource without putting.
	// It will return an exhausted error.
	resource, err = pool.Take(ctx)
	fmt.Println(resource, err, err == rego.ErrPoolExhausted)

	// Put the resource back to the pool.
	pool.Put(ctx, resource)
	pool.Close(ctx)

	// Now, the pool is closed so any taking from the pool will return a closed error.
	resource, err = pool.Take(ctx)
	fmt.Println(resource, err, err == rego.ErrPoolClosed)

	// Create a pool with limit and fast-failed and new error funcs.
	pool = rego.New(1, acquire, release, rego.WithFastFailed(), rego.WithPoolExhaustedErr(newPoolExhaustedErr), rego.WithPoolClosedErr(newPoolClosedErr))

	// Take one resource from pool which is ok.
	resource, err = pool.Take(ctx)
	fmt.Println(resource, err)

	// However, the pool is exhausted after taking one resource without putting.
	// It will return a customizing exhausted error.
	resource, err = pool.Take(ctx)
	fmt.Println(resource, err, err == errPoolExhausted)

	// Put the resource back to the pool.
	pool.Put(ctx, resource)
	pool.Close(ctx)

	// Now, the pool is closed so any taking from the pool will return a customizing closed error.
	resource, err = pool.Take(ctx)
	fmt.Println(resource, err, err == errPoolClosed)
}
