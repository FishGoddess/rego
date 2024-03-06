// Copyright 2024 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import "context"

type config struct {
	limit      uint64
	fastFailed bool

	newPoolFullErrFunc   func(ctx context.Context) error
	newPoolClosedErrFunc func(ctx context.Context) error
}

func newDefaultConfig() *config {
	conf := &config{
		limit:                64,
		fastFailed:           false,
		newPoolFullErrFunc:   nil,
		newPoolClosedErrFunc: nil,
	}

	return conf
}

type Option func(conf *config)

func (o Option) ApplyTo(conf *config) {
	o(conf)
}

// WithLimit sets limit to config.
func WithLimit(limit uint64) Option {
	return func(conf *config) {
		conf.limit = limit
	}
}

// WithFastFailed sets fastFailed to config.
func WithFastFailed() Option {
	return func(conf *config) {
		conf.fastFailed = true
	}
}

// WithPoolFullErr sets newPoolFullErr to config.
func WithPoolFullErr(newPoolFullErr func(ctx context.Context) error) Option {
	return func(conf *config) {
		conf.newPoolFullErrFunc = newPoolFullErr
	}
}

// WithPoolClosedErr sets newPoolClosedErr to config.
func WithPoolClosedErr(newPoolClosedErr func(ctx context.Context) error) Option {
	return func(conf *config) {
		conf.newPoolClosedErrFunc = newPoolClosedErr
	}
}
