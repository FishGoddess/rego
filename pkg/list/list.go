// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package list

import (
	stdlist "container/list"
	"fmt"
	"strings"
)

const (
	separator = "|"
)

type List[T any] struct {
	list *stdlist.List
}

func New[T any]() *List[T] {
	list := &List[T]{
		list: stdlist.New(),
	}

	return list
}

// Push pushes a value to list.
func (l *List[T]) Push(value T) {
	l.list.PushBack(value)
}

// Pop pops a value from list.
func (l *List[T]) Pop() (value T, ok bool) {
	front := l.list.Front()
	if front == nil {
		return value, false
	}

	value, ok = l.list.Remove(front).(T)
	return value, ok
}

// String stringifies a list with specified format.
func (l *List[T]) Remove(shouldRemove func(value T) bool) []T {
	var removedValues []T

	elem := l.list.Front()
	for elem != nil {
		current := elem
		elem = elem.Next()

		value, ok := current.Value.(T)
		if !ok {
			continue
		}

		if shouldRemove(value) {
			l.list.Remove(current)
			removedValues = append(removedValues, value)
		}
	}

	return removedValues
}

// Len returns the length of list.
func (l *List[T]) Len() uint64 {
	return uint64(l.list.Len())
}

// String stringifies a list with specified format.
func (l *List[T]) String() string {
	var builder strings.Builder

	elem := l.list.Front()
	for elem != nil {
		current := elem
		elem = elem.Next()

		fmt.Fprintf(&builder, "%v%s", current.Value, separator)
	}

	str := builder.String()
	str = strings.TrimSuffix(str, separator)
	return str
}
