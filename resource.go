// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

type resource[T any] struct {
	value T
}

func (r *resource[T]) reset() {
	var zero T
	r.value = zero
}
