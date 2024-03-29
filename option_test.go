// Copyright 2024 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
	"fmt"
	"testing"
)

// go test -v -cover -count=1 -test.cpu=1 -run=^TestWithLimit$
func TestWithLimit(t *testing.T) {
	conf := &config{limit: 0}
	WithLimit(64)(conf)

	if conf.limit != 64 {
		t.Fatalf("conf.limit %d is wrong", conf.limit)
	}
}

// go test -v -cover -count=1 -test.cpu=1 -run=^TestWithFastFailed$
func TestWithFastFailed(t *testing.T) {
	conf := &config{fastFailed: false}
	WithFastFailed()(conf)

	if !conf.fastFailed {
		t.Fatalf("conf.fastFailed %+v is wrong", conf.fastFailed)
	}
}

// go test -v -cover -count=1 -test.cpu=1 -run=^TestWithPoolFullErr$
func TestWithPoolFullErr(t *testing.T) {
	newPoolFullErr := func(ctx context.Context) error {
		return nil
	}

	conf := &config{newPoolFullErrFunc: nil}
	WithPoolFullErr(newPoolFullErr)(conf)

	if fmt.Sprintf("%p", conf.newPoolFullErrFunc) != fmt.Sprintf("%p", newPoolFullErr) {
		t.Fatalf("conf.newPoolFullErrFunc %p is wrong", conf.newPoolFullErrFunc)
	}
}

// go test -v -cover -count=1 -test.cpu=1 -run=^TestWithPoolClosedErr$
func TestWithPoolClosedErr(t *testing.T) {
	newPoolClosedErr := func(ctx context.Context) error {
		return nil
	}

	conf := &config{newPoolClosedErrFunc: nil}
	WithPoolClosedErr(newPoolClosedErr)(conf)

	if fmt.Sprintf("%p", conf.newPoolClosedErrFunc) != fmt.Sprintf("%p", newPoolClosedErr) {
		t.Fatalf("conf.newPoolClosedErrFunc %p is wrong", conf.newPoolClosedErrFunc)
	}
}
