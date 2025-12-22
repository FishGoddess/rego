// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import "testing"

// go test -v -cover -run=^TestResourceReset$
func TestResourceReset(t *testing.T) {
	got := resource[int]{
		value: 666,
	}

	got.reset()
	var want resource[int]
	if got != want {
		t.Fatalf("got %+v != want %+v", got, want)
	}
}
