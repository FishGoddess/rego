// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
	"errors"
)

var (
	ErrPoolExhausted = errors.New("rego: pool is exhausted")
	ErrPoolClosed    = errors.New("rego: pool is closed")
)

type config struct {
	newPoolExhaustedErr func(ctx context.Context) error
	newPoolClosedErr    func(ctx context.Context) error
}

func newConfig() *config {
	newPoolExhaustedErr := func(_ context.Context) error {
		return ErrPoolExhausted
	}

	newPoolClosedErr := func(_ context.Context) error {
		return ErrPoolClosed
	}

	conf := &config{
		newPoolExhaustedErr: newPoolExhaustedErr,
		newPoolClosedErr:    newPoolClosedErr,
	}

	return conf
}

func (c *config) apply(opts ...Option) *config {
	for _, opt := range opts {
		opt(c)
	}

	return c
}

type Option func(conf *config)

// WithPoolExhaustedErr sets a function returns an error of exhausted pool to config.
func WithPoolExhaustedErr(newErr func(ctx context.Context) error) Option {
	return func(conf *config) {
		if newErr != nil {
			conf.newPoolExhaustedErr = newErr
		}
	}
}

// WithPoolClosedErr sets a function returns an error of closed pool to config.
func WithPoolClosedErr(newErr func(ctx context.Context) error) Option {
	return func(conf *config) {
		if newErr != nil {
			conf.newPoolClosedErr = newErr
		}
	}
}
