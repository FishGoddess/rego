// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package list

import (
	"container/list"
	"testing"
	"time"
)

// go test -v -cover -run=^TestNewList$
func TestNewList(t *testing.T) {
	nodeList := list.New()

	list := &List[int]{
		nodeList: nodeList,
	}

	gotString := list.String()
	wantString := ""

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

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

	nodeList.PushBack(node1)
	nodeList.PushBack(node2)
	nodeList.PushBack(node3)

	gotString = list.String()
	wantString = "1|2025-05-03 22:17:08|2026-11-15 12:08:53|2|1997-12-23 02:01:45|2008-11-23 12:10:01|3|1999-01-02 03:05:06|2012-09-19 12:08:10"

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}
}

// go test -v -cover -run=^TestList$
func TestList(t *testing.T) {
	defer func() {
		nowFunc = time.Now
	}()

	nodes := make([]*Node[int], 0, 3)
	list := NewList[int]()
	length := uint64(0)

	nowFunc = func() time.Time {
		return time.Date(2025, 5, 3, 22, 17, 8, 0, time.Local)
	}

	node := NewNode(1)
	nodes = append(nodes, node)
	list.Push(node)
	length++

	nowFunc = func() time.Time {
		return time.Date(1997, 12, 23, 2, 1, 45, 0, time.Local)
	}

	node = NewNode(2)
	nodes = append(nodes, node)
	list.Push(node)
	length++

	nowFunc = func() time.Time {
		return time.Date(1999, 1, 2, 3, 5, 6, 0, time.Local)
	}

	node = NewNode(3)
	nodes = append(nodes, node)
	list.Push(node)
	length++

	gotString := list.String()
	wantString := "1|2025-05-03 22:17:08|2025-05-03 22:17:08|2|1997-12-23 02:01:45|1997-12-23 02:01:45|3|1999-01-02 03:05:06|1999-01-02 03:05:06"

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != length {
		t.Fatalf("list.Len() %d != length %d", list.Len(), length)
	}

	if list.Len() != uint64(len(nodes)) {
		t.Fatalf("list.Len() %d != len(nodes) %d", list.Len(), len(nodes))
	}

	for i := range list.Len() {
		node, ok := list.Pop()
		if !ok {
			t.Fatal("list pop not ok")
		}

		if node != nodes[i] {
			t.Fatalf("node %s != nodes[i] %s", node, nodes[i])
		}

		length--

		if list.Len() != length {
			t.Fatalf("list.Len() %d != length %d", list.Len(), length)
		}
	}

	if _, ok := list.Pop(); ok {
		t.Fatalf("list pop sok")
	}
}

// go test -v -cover -run=^TestListRemove$
func TestListRemove(t *testing.T) {
	list := NewList[int]()
	length := uint64(0)

	nowFunc = func() time.Time {
		return time.Date(2025, 5, 3, 22, 17, 8, 0, time.Local)
	}

	node := NewNode(1)
	list.Push(node)
	length++

	nowFunc = func() time.Time {
		return time.Date(1997, 12, 23, 2, 1, 45, 0, time.Local)
	}

	node = NewNode(2)
	list.Push(node)
	length++

	nowFunc = func() time.Time {
		return time.Date(1999, 1, 2, 3, 5, 6, 0, time.Local)
	}

	node = NewNode(3)
	list.Push(node)
	length++

	gotString := list.String()
	wantString := "1|2025-05-03 22:17:08|2025-05-03 22:17:08|2|1997-12-23 02:01:45|1997-12-23 02:01:45|3|1999-01-02 03:05:06|1999-01-02 03:05:06"

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != length {
		t.Fatalf("list.Len() %d != length %d", list.Len(), length)
	}

	createTime := time.Date(1997, 12, 23, 2, 1, 44, 0, time.Local)
	list.Remove(createTime)

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != length {
		t.Fatalf("list.Len() %d != length %d", list.Len(), length)
	}

	createTime = time.Date(1997, 12, 23, 2, 1, 45, 0, time.Local)
	list.Remove(createTime)
	length--

	gotString = list.String()
	wantString = "1|2025-05-03 22:17:08|2025-05-03 22:17:08|3|1999-01-02 03:05:06|1999-01-02 03:05:06"

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != length {
		t.Fatalf("list.Len() %d != length %d", list.Len(), length)
	}

	createTime = time.Date(1999, 1, 2, 3, 5, 6, 0, time.Local)
	list.Remove(createTime)
	length--

	gotString = list.String()
	wantString = "1|2025-05-03 22:17:08|2025-05-03 22:17:08"

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != length {
		t.Fatalf("list.Len() %d != length %d", list.Len(), length)
	}

	createTime = time.Date(2025, 5, 3, 22, 17, 8, 0, time.Local)
	list.Remove(createTime)
	length--

	gotString = list.String()
	wantString = ""

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != length {
		t.Fatalf("list.Len() %d != length %d", list.Len(), length)
	}
}
