// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package list

import (
	"container/list"
	"strings"
	"time"
)

type List[T any] struct {
	nodeList *list.List
}

func NewList[T any]() *List[T] {
	list := &List[T]{
		nodeList: list.New(),
	}

	return list
}

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

func (l *List[T]) Push(node *Node[T]) {
	l.nodeList.PushBack(node)
}

func (l *List[T]) Pop() (node *Node[T], ok bool) {
	front := l.nodeList.Front()
	if front == nil {
		return nil, false
	}

	node, ok = l.nodeList.Remove(front).(*Node[T])
	return node, ok
}

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

func (l *List[T]) Len() int {
	return l.nodeList.Len()
}
