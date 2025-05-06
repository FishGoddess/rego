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
	disableToken bool

	newPoolExhaustedErrFunc func(ctx context.Context) error
	newPoolClosedErrFunc    func(ctx context.Context) error
}

func newDefaultConfig() *config {
	newPoolExhaustedErr := func(_ context.Context) error {
		return ErrPoolExhausted
	}

	newPoolClosedErr := func(_ context.Context) error {
		return ErrPoolClosed
	}

	conf := &config{
		disableToken:            false,
		newPoolExhaustedErrFunc: newPoolExhaustedErr,
		newPoolClosedErrFunc:    newPoolClosedErr,
	}

	return conf
}

func (c *config) newPoolExhaustedErr(ctx context.Context) error {
	newPoolExhaustedErr := c.newPoolExhaustedErrFunc
	return newPoolExhaustedErr(ctx)
}

func (c *config) newPoolClosedErr(ctx context.Context) error {
	newPoolClosedErr := c.newPoolClosedErrFunc
	return newPoolClosedErr(ctx)
}

type Option func(conf *config)

func (o Option) ApplyTo(conf *config) {
	o(conf)
}

// WithDisableToken sets disableToken to config.
func WithDisableToken() Option {
	return func(conf *config) {
		conf.disableToken = true
	}
}

// WithPoolExhaustedErr sets newPoolExhaustedErr to config.
func WithPoolExhaustedErr(newPoolExhaustedErr func(ctx context.Context) error) Option {
	return func(conf *config) {
		if newPoolExhaustedErr != nil {
			conf.newPoolExhaustedErrFunc = newPoolExhaustedErr
		}
	}
}

// WithPoolClosedErr sets newPoolClosedErr to config.
func WithPoolClosedErr(newPoolClosedErr func(ctx context.Context) error) Option {
	return func(conf *config) {
		if newPoolClosedErr != nil {
			conf.newPoolClosedErrFunc = newPoolClosedErr
		}
	}
}
