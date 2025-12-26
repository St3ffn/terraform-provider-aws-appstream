// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"time"
)

const (
	defaultRetryTimeout     = 60 * time.Minute
	defaultRetryInitBackoff = 1 * time.Second
	defaultRetryMaxBackoff  = 5 * time.Minute
)

type retryOnFn func(error) bool

type retryConfig struct {
	timeout     time.Duration
	initBackoff time.Duration
	maxBackoff  time.Duration
	retryOnFns  []retryOnFn
}

func defaultRetryConfig() *retryConfig {
	return &retryConfig{
		timeout:     defaultRetryTimeout,
		initBackoff: defaultRetryInitBackoff,
		maxBackoff:  defaultRetryMaxBackoff,
		retryOnFns:  []retryOnFn{},
	}
}

type retryOption func(*retryConfig)

//nolint:unused // reserved for future retry configuration options
func withTimeout(timeout time.Duration) retryOption {
	return func(c *retryConfig) {
		c.timeout = timeout
	}
}

//nolint:unused // reserved for future retry configuration options
func withInitBackoff(initBackoff time.Duration) retryOption {
	return func(c *retryConfig) {
		c.initBackoff = initBackoff
	}
}

//nolint:unused // reserved for future retry configuration options
func withMaxBackoff(maxBackoff time.Duration) retryOption {
	return func(c *retryConfig) {
		c.maxBackoff = maxBackoff
	}
}

func withRetryOnFns(retryOnFns ...retryOnFn) retryOption {
	return func(c *retryConfig) {
		c.retryOnFns = append(c.retryOnFns, retryOnFns...)
	}
}

func shouldRetry(err error, retryOnFns []retryOnFn) bool {
	for _, retryFn := range retryOnFns {
		if retryFn(err) {
			return true
		}
	}
	return false
}

func retryOn(ctx context.Context, call func(context.Context) error, opts ...retryOption) error {
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
