// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package list

import (
	"container/list"
	"strings"
	"sync"
	"time"
)

var (
	nowFunc = time.Now
)

type List[T any] struct {
	nodeList *list.List
	nodePool *sync.Pool
}

func New[T any]() *List[T] {
	nodeList := list.New()

	nodePool := &sync.Pool{
		New: func() any {
			return new(Node[T])
		},
	}

	list := &List[T]{
		nodeList: nodeList,
		nodePool: nodePool,
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

func (l *List[T]) newNode(value T) *Node[T] {
	now := nowFunc()

	node, ok := l.nodePool.Get().(*Node[T])
	if !ok {
		node = new(Node[T])
	}

	node.value = value
	node.createTime = now
	node.updateTime = now
	return node
}

func (l *List[T]) freeNode(node *Node[T]) {
	var value T
	node.value = value

	l.nodePool.Put(node)
}

func (l *List[T]) Push(value T) {
	node := l.newNode(value)
	l.nodeList.PushBack(node)
}

func (l *List[T]) Pop() (value T, ok bool) {
	front := l.nodeList.Front()
	if front == nil {
		return value, false
	}

	node, ok := l.nodeList.Remove(front).(*Node[T])
	if !ok {
		return value, false
	}

	value = node.value
	l.freeNode(node)

	return value, true
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

		node, ok = l.nodeList.Remove(current).(*Node[T])
		if !ok {
			continue
		}

		l.freeNode(node)
	}
}

func (l *List[T]) Len() int {
	return l.nodeList.Len()
}
