// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package list

import (
	"testing"
	"time"
)

// go test -v -cover -run=^TestNodeString$
func TestNodeString(t *testing.T) {
	node := &Node[int]{
		value:      1,
		createTime: time.Date(2025, 5, 3, 22, 17, 8, 0, time.Local),
		updateTime: time.Date(2026, 11, 15, 12, 8, 53, 0, time.Local),
	}

	gotString := node.String()
	wantString := "1|2025-05-03 22:17:08|2026-11-15 12:08:53"

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}
}
