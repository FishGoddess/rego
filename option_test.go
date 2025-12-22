// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
	"fmt"
	"testing"
)

// go test -v -cover -run=^TestWithPoolExhaustedErr$
func TestWithPoolExhaustedErr(t *testing.T) {
	newErr := func(ctx context.Context) error {
		return nil
	}

	conf := &config{newPoolExhaustedErr: nil}
	WithPoolExhaustedErr(newErr)(conf)

	got := fmt.Sprintf("%p", conf.newPoolExhaustedErr)
	want := fmt.Sprintf("%p", newErr)
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}

	WithPoolExhaustedErr(nil)(conf)

	got = fmt.Sprintf("%p", conf.newPoolExhaustedErr)
	want = fmt.Sprintf("%p", newErr)
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}

// go test -v -cover -run=^TestWithPoolClosedErr$
func TestWithPoolClosedErr(t *testing.T) {
	newErr := func(ctx context.Context) error {
		return nil
	}

	conf := &config{newPoolClosedErr: nil}
	WithPoolClosedErr(newErr)(conf)

	got := fmt.Sprintf("%p", conf.newPoolClosedErr)
	want := fmt.Sprintf("%p", newErr)
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}

	WithPoolClosedErr(nil)(conf)

	got = fmt.Sprintf("%p", conf.newPoolClosedErr)
	want = fmt.Sprintf("%p", newErr)
	if got != want {
		t.Fatalf("got %s != want %s", got, want)
	}
}
