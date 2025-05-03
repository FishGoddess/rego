// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package list

import (
	"container/list"
	"testing"
	"time"
)

// go test -v -cover -run=^TestListString$
func TestListString(t *testing.T) {
	node1 := &Node[int]{
		value:      1,
		createTime: time.Date(2025, 5, 3, 22, 17, 8, 0, time.Local),
		updateTime: time.Date(2026, 11, 15, 12, 8, 53, 0, time.Local),
	}

	node2 := &Node[int]{
		value:      2,
		createTime: time.Date(1997, 12, 23, 2, 1, 45, 0, time.Local),
		updateTime: time.Date(2008, 11, 23, 12, 10, 1, 0, time.Local),
	}

	node3 := &Node[int]{
		value:      3,
		createTime: time.Date(1999, 1, 2, 3, 5, 6, 0, time.Local),
		updateTime: time.Date(2012, 9, 19, 12, 8, 10, 0, time.Local),
	}

	nodeList := list.New()
	nodeList.PushBack(node1)
	nodeList.PushBack(node2)
	nodeList.PushBack(node3)

	list := &List[int]{
		nodeList: nodeList,
	}

	gotString := list.String()
	wantString := "1|2025-05-03 22:17:08|2026-11-15 12:08:53|2|1997-12-23 02:01:45|2008-11-23 12:10:01|3|1999-01-02 03:05:06|2012-09-19 12:08:10"

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}
}

// go test -v -cover -run=^TestListString$
func TestList(t *testing.T) {
	node1 := &Node[int]{
		value:      1,
		createTime: time.Date(2025, 5, 3, 22, 17, 8, 0, time.Local),
		updateTime: time.Date(2026, 11, 15, 12, 8, 53, 0, time.Local),
	}

	node2 := &Node[int]{
		value:      2,
		createTime: time.Date(1997, 12, 23, 2, 1, 45, 0, time.Local),
		updateTime: time.Date(2008, 11, 23, 12, 10, 1, 0, time.Local),
	}

	node3 := &Node[int]{
		value:      3,
		createTime: time.Date(1999, 1, 2, 3, 5, 6, 0, time.Local),
		updateTime: time.Date(2012, 9, 19, 12, 8, 10, 0, time.Local),
	}

	nodeList := list.New()
	nodeList.PushBack(node1)
	nodeList.PushBack(node2)
	nodeList.PushBack(node3)

	list := &List[int]{
		nodeList: nodeList,
	}

	gotString := list.String()
	wantString := "1|2025-05-03 22:17:08|2026-11-15 12:08:53|2|1997-12-23 02:01:45|2008-11-23 12:10:01|3|1999-01-02 03:05:06|2012-09-19 12:08:10"

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}
}
