// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import "time"

// Status includes some statistics of the pool.
type Status struct {
	// Limit is the maximum quantity of resources in pool.
	Limit uint64 `json:"limit"`

	// Using is the quantity of using resources in pool.
	Using uint64 `json:"using"`

	// Idle is the quantity of idle resources in pool.
	Idle uint64 `json:"idle"`

	// Waiting is the quantity of caller waiting for a resource.
	Waiting uint64 `json:"waiting"`

	// WaitDuration is the average duration waiting a resource.
	WaitDuration time.Duration `json:"wait_duration"`
}
