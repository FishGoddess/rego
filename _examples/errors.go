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
	errPoolIsFull   = errors.New("_examples: pool is full")
	errPoolIsClosed = errors.New("_examples: pool is closed")
)

func acquire() (int, error) {
	return 0, nil
}

func release(resource int) error {
	return nil
}

func newPoolFullErr(ctx context.Context) error {
	return errPoolIsFull
}

func newPoolClosedErr(ctx context.Context) error {
	return errPoolIsClosed
}

func main() {
	// Create a pool with limit and fast-failed so it will return an error immediately instead of waiting.
	ctx := context.Background()
	pool := rego.New(1, acquire, release, rego.WithFastFailed())

	// Take one resource from pool which is ok.
	resource, err := pool.Take(ctx)
	fmt.Println(resource, err)

	// However, the pool is full after taking one resource without putting.
	// It will return a full error.
	resource, err = pool.Take(ctx)
	fmt.Println(resource, err, err == rego.ErrPoolIsFull)

	// Put the resource back to the pool.
	pool.Put(resource)
	pool.Close()

	// Now, the pool is closed so any taking from the pool will return a closed error.
	resource, err = pool.Take(ctx)
	fmt.Println(resource, err, err == rego.ErrPoolIsClosed)

	// Create a pool with limit and fast-failed and new error funcs.
	pool = rego.New(1, acquire, release, rego.WithFastFailed(), rego.WithPoolFullErr(newPoolFullErr), rego.WithPoolClosedErr(newPoolClosedErr))

	// Take one resource from pool which is ok.
	resource, err = pool.Take(ctx)
	fmt.Println(resource, err)

	// However, the pool is full after taking one resource without putting.
	// It will return a customizing full error.
	resource, err = pool.Take(ctx)
	fmt.Println(resource, err, err == errPoolIsFull)

	// Put the resource back to the pool.
	pool.Put(resource)
	pool.Close()

	// Now, the pool is closed so any taking from the pool will return a customizing closed error.
	resource, err = pool.Take(ctx)
	fmt.Println(resource, err, err == errPoolIsClosed)
}
