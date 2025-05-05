// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package list

import (
	"container/list"
	"strings"
	"time"
)

// List stores some nodes with time control mechanism.
type List[T any] struct {
	nodeList *list.List
}

// NewList creates an empty list.
func NewList[T any]() *List[T] {
	list := &List[T]{
		nodeList: list.New(),
	}

	return list
}

// String stringifies a list with specified format.
func (l *List[T]) String() string {
	var builder strings.Builder

	elem := l.nodeList.Front()
	for elem != nil {
		current := elem
		elem = elem.Next()

		node, ok := current.Value.(*Node[T])
		if !ok {
			continue
		}

		builder.WriteString(node.String())

		if elem != nil {
			builder.WriteByte(fieldSeparator)
		}
	}

	return builder.String()
}

// Push pushes a node to list.
func (l *List[T]) Push(node *Node[T]) {
	l.nodeList.PushBack(node)
}

// Pop pops a node from list.
func (l *List[T]) Pop() (node *Node[T], ok bool) {
	front := l.nodeList.Front()
	if front == nil {
		return nil, false
	}

	node, ok = l.nodeList.Remove(front).(*Node[T])
	return node, ok
}

// Remove removes the nodes whose create time is earlier than the given create time.
func (l *List[T]) Remove(createTime time.Time) {
	elem := l.nodeList.Front()
	for elem != nil {
		current := elem
		elem = elem.Next()

		node, ok := current.Value.(*Node[T])
		if !ok {
			continue
		}

		if node.createTime.After(createTime) {
			continue
		}

		l.nodeList.Remove(current)
	}
}

// Len returns the length of list.
func (l *List[T]) Len() uint64 {
	return uint64(l.nodeList.Len())
}
