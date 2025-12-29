// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"time"
)

const (
	defaultRetryTimeout     = 60 * time.Minute
	defaultRetryInitBackoff = 1 * time.Second
	defaultRetryMaxBackoff  = 5 * time.Minute
)

type RetryOnFn func(error) bool

type retryConfig struct {
	timeout     time.Duration
	initBackoff time.Duration
	maxBackoff  time.Duration
	retryOnFns  []RetryOnFn
}

func defaultRetryConfig() *retryConfig {
	return &retryConfig{
		timeout:     defaultRetryTimeout,
		initBackoff: defaultRetryInitBackoff,
		maxBackoff:  defaultRetryMaxBackoff,
		retryOnFns:  []RetryOnFn{},
	}
}

type RetryOption func(*retryConfig)

//nolint:unused // reserved for future retry configuration options
func WithTimeout(timeout time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.timeout = timeout
	}
}

//nolint:unused // reserved for future retry configuration options
func WithInitBackoff(initBackoff time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.initBackoff = initBackoff
	}
}

//nolint:unused // reserved for future retry configuration options
func WithMaxBackoff(maxBackoff time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.maxBackoff = maxBackoff
	}
}

func WithRetryOnFns(retryOnFns ...RetryOnFn) RetryOption {
	return func(c *retryConfig) {
		c.retryOnFns = append(c.retryOnFns, retryOnFns...)
	}
}

func shouldRetry(err error, retryOnFns []RetryOnFn) bool {
	for _, retryFn := range retryOnFns {
		if retryFn(err) {
			return true
		}
	}
	return false
}

func RetryOn(ctx context.Context, call func(context.Context) error, opts ...RetryOption) error {
	options := defaultRetryConfig()
	for _, fn := range opts {
		fn(options)
	}

	ctx, cancel := context.WithTimeout(ctx, options.timeout)
	defer cancel()

	backoff := options.initBackoff

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := call(ctx)
		if err == nil {
			return nil
		}

		if !shouldRetry(err, options.retryOnFns) {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}

		backoff *= 2
		if backoff > options.maxBackoff {
			backoff = options.maxBackoff
		}
	}
}
