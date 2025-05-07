// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package list

import (
	"container/list"
	"testing"
)

// go test -v -run=^$ -bench=^BenchmarkContainerList$ -benchtime=1s
func BenchmarkContainerList(b *testing.B) {
	list := list.New()

	b.ReportAllocs()
	b.ResetTimer()

	for i := range b.N {
		elem := list.PushBack(i)
		list.Remove(elem)
	}
}

// go test -v -run=^$ -bench=^BenchmarkList$ -benchtime=1s
func BenchmarkList(b *testing.B) {
	list := New[int]()

	b.ReportAllocs()
	b.ResetTimer()

	for i := range b.N {
		list.Push(i)
		list.Pop()
	}
}
