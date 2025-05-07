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
	testCases := [][]int{
		{},
		{1},
		{2},
		{1, 2},
		{2, 1},
		{1, 3, 5},
		{2, 4, 6},
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	for _, values := range testCases {
		list := New[int]()
		length := uint64(len(values))

		strs := make([]string, 0, len(values))
		for _, value := range values {
			list.Push(value)
			strs = append(strs, strconv.Itoa(value))
		}

		gotString := list.String()
		wantString := strings.Join(strs, separator)

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
}

// go test -v -cover -run=^TestListRemove$
func TestListRemove(t *testing.T) {
	testCases := [][]int{
		{},
		{1},
		{2},
		{1, 2},
		{2, 1},
		{1, 3, 5},
		{2, 4, 6},
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	for _, values := range testCases {
		list := New[int]()
		length := uint64(len(values))

		strs := make([]string, 0, len(values))
		for _, value := range values {
			list.Push(value)
			strs = append(strs, strconv.Itoa(value))
		}

		gotString := list.String()
		wantString := strings.Join(strs, separator)

		if gotString != wantString {
			t.Fatalf("got string %s != want string %s", gotString, wantString)
		}

		if list.Len() != length {
			t.Fatalf("list.Len() %d != length %d", list.Len(), length)
		}

		shouldRemove := func(value int) bool {
			return value%2 != 0
		}

		wantRemovedValues := make([]int, 0, len(values))
		for _, value := range values {
			if shouldRemove(value) {
				wantRemovedValues = append(wantRemovedValues, value)
			}
		}

		removedValues := list.Remove(shouldRemove)

		if !slices.Equal(removedValues, wantRemovedValues) {
			t.Fatalf("removedValues %v != wantRemovedValues %v", removedValues, wantRemovedValues)
		}

		strs = make([]string, 0, len(removedValues))
		for _, value := range values {
			if !shouldRemove(value) {
				strs = append(strs, strconv.Itoa(value))
			}
		}

		gotString = list.String()
		wantString = strings.Join(strs, separator)

		if gotString != wantString {
			t.Fatalf("got string %s != want string %s", gotString, wantString)
		}

		wantLen := length - uint64(len(removedValues))
		if list.Len() != wantLen {
			t.Fatalf("list.Len() %d != wantLen %d", list.Len(), wantLen)
		}
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
