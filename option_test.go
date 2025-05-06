// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

// go test -v -cover -run=^TestWithFastFailed$
func TestWithFastFailed(t *testing.T) {
	conf := &config{fastFailed: false}
	WithFastFailed()(conf)

	if !conf.fastFailed {
		t.Fatalf("conf.fastFailed %+v is wrong", conf.fastFailed)
	}
}

// go test -v -cover -run=^TestWithPoolExhaustedErr$
func TestWithPoolExhaustedErr(t *testing.T) {
	errPoolExhausted := errors.New("exhausted")

	newPoolExhaustedErr := func(ctx context.Context) error {
		return errPoolExhausted
	}

	conf := &config{newPoolExhaustedErrFunc: nil}
	WithPoolExhaustedErr(newPoolExhaustedErr)(conf)

	if fmt.Sprintf("%p", conf.newPoolExhaustedErrFunc) != fmt.Sprintf("%p", newPoolExhaustedErr) {
		t.Fatalf("conf.newPoolExhaustedErr %p is wrong", conf.newPoolExhaustedErrFunc)
	}

	ctx := context.Background()
	if err := conf.newPoolExhaustedErr(ctx); err != errPoolExhausted {
		t.Fatalf("err %v != errPoolExhausted %v", err, errPoolExhausted)
	}

	WithPoolExhaustedErr(nil)(conf)

	if fmt.Sprintf("%p", conf.newPoolExhaustedErrFunc) != fmt.Sprintf("%p", newPoolExhaustedErr) {
		t.Fatalf("conf.newPoolExhaustedErr %p is wrong", conf.newPoolExhaustedErrFunc)
	}
}

// go test -v -cover -run=^TestWithPoolClosedErr$
func TestWithPoolClosedErr(t *testing.T) {
	errPoolClosed := errors.New("closed")

	newPoolClosedErr := func(ctx context.Context) error {
		return errPoolClosed
	}

	conf := &config{newPoolClosedErrFunc: nil}
	WithPoolClosedErr(newPoolClosedErr)(conf)

	if fmt.Sprintf("%p", conf.newPoolClosedErrFunc) != fmt.Sprintf("%p", newPoolClosedErr) {
		t.Fatalf("conf.newPoolClosedErr %p is wrong", conf.newPoolClosedErrFunc)
	}

	ctx := context.Background()
	if err := conf.newPoolClosedErr(ctx); err != errPoolClosed {
		t.Fatalf("err %v != errPoolClosed %v", err, errPoolClosed)
	}

	WithPoolClosedErr(nil)(conf)

	if fmt.Sprintf("%p", conf.newPoolClosedErrFunc) != fmt.Sprintf("%p", newPoolClosedErr) {
		t.Fatalf("conf.newPoolClosedErr %p is wrong", conf.newPoolClosedErrFunc)
	}
}
