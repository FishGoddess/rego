// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rego

import (
	"context"
	"errors"
)

var (
	ErrPoolIsFull   = errors.New("rego: pool is full")
	ErrPoolIsClosed = errors.New("rego: pool is closed")
)

type config struct {
	fastFailed bool

	newPoolFullErr   func(ctx context.Context) error
	newPoolClosedErr func(ctx context.Context) error
}

func newDefaultConfig() *config {
	newPoolFullErr := func(_ context.Context) error {
		return ErrPoolIsFull
	}

	newPoolClosedErr := func(_ context.Context) error {
		return ErrPoolIsClosed
	}

	conf := &config{
		fastFailed:       false,
		newPoolFullErr:   newPoolFullErr,
		newPoolClosedErr: newPoolClosedErr,
	}

	return conf
}

type Option func(conf *config)

func (o Option) ApplyTo(conf *config) {
	o(conf)
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
		if newPoolFullErr != nil {
			conf.newPoolFullErr = newPoolFullErr
		}
	}
}

// WithPoolClosedErr sets newPoolClosedErr to config.
func WithPoolClosedErr(newPoolClosedErr func(ctx context.Context) error) Option {
	return func(conf *config) {
		if newPoolClosedErr != nil {
			conf.newPoolClosedErr = newPoolClosedErr
		}
	}
}
