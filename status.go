// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import "time"

// PoolStatus is the statistics of pool.
type PoolStatus struct {
	// Limit is the maximum quantity of resources in pool.
	Limit uint64 `json:"limit"`

	// Active is the quantity of resources in pool including idle and using.
	Active uint64 `json:"active"`

	// Idle is the quantity of idle resources in pool.
	Idle uint64 `json:"idle"`

	// Waiting is the quantity of waiting for a resource.
	Waiting uint64 `json:"waiting"`

	// AverageWaitDuration is the average wait duration for new resources.
	AverageWaitDuration time.Duration `json:"average_wait_duration"`
}
