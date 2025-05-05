// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package list

import (
	"fmt"
	"time"
)

const (
	fieldSeparator = '|'
	timeFormat     = "2006-01-02 15:04:05"
)

var (
	nowFunc = time.Now
)

type Node[T any] struct {
	value T

	createTime time.Time
	updateTime time.Time
}

func NewNode[T any](value T) *Node[T] {
	now := nowFunc()

	node := &Node[T]{
		value:      value,
		createTime: now,
		updateTime: now,
	}

	return node
}

func (n *Node[T]) String() string {
	bs := make([]byte, 0, 64)

	bs = fmt.Append(bs, n.value)
	bs = append(bs, fieldSeparator)
	bs = n.createTime.AppendFormat(bs, timeFormat)
	bs = append(bs, fieldSeparator)
	bs = n.updateTime.AppendFormat(bs, timeFormat)
	return string(bs)
}
