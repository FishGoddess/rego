// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package list

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

// go test -v -cover -run=^TestList$
func TestList(t *testing.T) {
	list := New[int]()

	gotString := list.String()
	wantString := ""

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != 0 {
		t.Fatalf("list.Len() %d is wrong", list.Len())
	}

	values := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	length := uint64(len(values))

	strs := make([]string, 0, len(values))
	for _, value := range values {
		list.Push(value)
		strs = append(strs, strconv.Itoa(value))
	}

	gotString = list.String()
	wantString = strings.Join(strs, separator)

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != length {
		t.Fatalf("list.Len() %d != length %d", list.Len(), length)
	}

	for _, value := range values {
		popValue, ok := list.Pop()
		if !ok {
			t.Fatal("list pop not ok")
		}

		if popValue != value {
			t.Fatalf("popValue %d != value %d", popValue, value)
		}

		length--

		if list.Len() != length {
			t.Fatalf("list.Len() %d != length %d", list.Len(), length)
		}
	}

	if _, ok := list.Pop(); ok {
		t.Fatal("list pop is ok")
	}

	if list.Len() != 0 {
		t.Fatalf("list.Len() %d is wrong", list.Len())
	}

	length = uint64(len(values))

	strs = make([]string, 0, len(values))
	for _, value := range values {
		list.Push(value)
		strs = append(strs, strconv.Itoa(value))
	}

	gotString = list.String()
	wantString = strings.Join(strs, separator)

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != length {
		t.Fatalf("list.Len() %d != length %d", list.Len(), length)
	}
}

// go test -v -cover -run=^TestListRemove$
func TestListRemove(t *testing.T) {
	list := New[int]()

	list.Push(1)
	list.Push(2)
	list.Push(3)
	list.Push(4)
	list.Push(5)

	gotString := list.String()
	wantString := "1|2|3|4|5"

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != 5 {
		t.Fatalf("list.Len() %d is wrong", list.Len())
	}

	shouldRemove := func(value int) bool {
		return value%2 == 0
	}

	removedValues := list.Remove(shouldRemove)
	wantValues := []int{2, 4}

	if !slices.Equal(removedValues, wantValues) {
		t.Fatalf("removedValues %v != wantValues %v", removedValues, wantValues)
	}

	gotString = list.String()
	wantString = "1|3|5"

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != 3 {
		t.Fatalf("list.Len() %d is wrong", list.Len())
	}
}

// go test -v -cover -run=^TestListString$
func TestListString(t *testing.T) {
	list := New[int]()

	gotString := list.String()
	wantString := ""

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != 0 {
		t.Fatalf("list.Len() %d is wrong", list.Len())
	}

	list.Push(1)
	list.Push(2)
	list.Push(3)

	gotString = list.String()
	wantString = "1|2|3"

	if gotString != wantString {
		t.Fatalf("got string %s != want string %s", gotString, wantString)
	}

	if list.Len() != 3 {
		t.Fatalf("list.Len() %d is wrong", list.Len())
	}
}
