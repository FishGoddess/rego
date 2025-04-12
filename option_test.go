// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
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

// go test -v -cover -run=^TestWithPoolFullErr$
func TestWithPoolFullErr(t *testing.T) {
	newPoolFullErr := func(ctx context.Context) error {
		return nil
	}

	conf := &config{newPoolFullErrFunc: nil}
	WithPoolFullErr(newPoolFullErr)(conf)

	if fmt.Sprintf("%p", conf.newPoolFullErrFunc) != fmt.Sprintf("%p", newPoolFullErr) {
		t.Fatalf("conf.newPoolFullErr %p is wrong", conf.newPoolFullErrFunc)
	}
}

// go test -v -cover -run=^TestWithPoolClosedErr$
func TestWithPoolClosedErr(t *testing.T) {
	newPoolClosedErr := func(ctx context.Context) error {
		return nil
	}

	conf := &config{newPoolClosedErrFunc: nil}
	WithPoolClosedErr(newPoolClosedErr)(conf)

	if fmt.Sprintf("%p", conf.newPoolClosedErrFunc) != fmt.Sprintf("%p", newPoolClosedErr) {
		t.Fatalf("conf.newPoolClosedErr %p is wrong", conf.newPoolClosedErrFunc)
	}
}
