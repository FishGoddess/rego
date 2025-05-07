// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package list

import (
	"fmt"
	"strings"
	"sync"
)

const (
	separator = "|"
)

type Element[T any] struct {
	value      T
	prev, next *Element[T]
}

type List[T any] struct {
	head *Element[T]
	tail *Element[T]
	len  uint64

	elementPool *sync.Pool
}

func New[T any]() *List[T] {
	elementPool := &sync.Pool{
		New: func() any {
			return new(Element[T])
		},
	}

	list := &List[T]{
		head:        nil,
		tail:        nil,
		len:         0,
		elementPool: elementPool,
	}

	return list
}

func (l *List[T]) newElement(value T) *Element[T] {
	elem := l.elementPool.Get().(*Element[T])
	elem.value = value

	return elem
}

func (l *List[T]) freeElement(elem *Element[T]) T {
	value := elem.value

	var zeroValue T
	elem.value = zeroValue
	elem.prev = nil
	elem.next = nil

	l.elementPool.Put(elem)
	return value
}

func (l *List[T]) Push(value T) {
	elem := l.newElement(value)

	if l.len == 0 {
		l.head = elem
	} else {
		elem.prev = l.tail
		l.tail.next = elem
	}

	l.tail = elem
	l.len++
}

func (l *List[T]) Pop() (value T, ok bool) {
	if l.len == 0 {
		return value, false
	}

	elem := l.head

	// Just need setting head and tail to nil if there is only one element in list.
	// This 'if else' is more readable than not distinguishing length.
	if l.len == 1 {
		l.head = nil
		l.tail = nil
	} else {
		l.head = l.head.next
		l.head.prev = nil
	}

	l.len--

	value = l.freeElement(elem)
	return value, true
}

func (l *List[T]) Remove(shouldRemove func(value T) bool) []T {
	var removedValues []T

	elem := l.head
	for elem != nil {
		current := elem
		elem = elem.next

		if shouldRemove(current.value) {
			if current == l.head {
				l.head = current.next
			}

			if current == l.tail {
				l.tail = current.prev
			}

			if current.prev != nil {
				current.prev.next = current.next
			}

			if current.next != nil {
				current.next.prev = current.prev
			}

			l.len--

			value := l.freeElement(current)
			removedValues = append(removedValues, value)
		}
	}

	return removedValues
}

func (l *List[T]) Len() uint64 {
	return l.len
}

func (l *List[T]) String() string {
	var builder strings.Builder

	elem := l.head
	for elem != nil {
		fmt.Fprintf(&builder, "%v%s", elem.value, separator)
		elem = elem.next
	}

	str := builder.String()
	str = strings.TrimSuffix(str, separator)
	return str
}
