// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"errors"
	"time"
)

const (
	defaultRetryTimeout     = 20 * time.Minute
	defaultRetryInitBackoff = 5 * time.Second
	defaultRetryMaxBackoff  = 1 * time.Minute
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

func WithTimeout(timeout time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.timeout = timeout
	}
}

func WithInitBackoff(initBackoff time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.initBackoff = initBackoff
	}
}

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

func contextTerminationError(ctxErr, lastErr error) error {
	if ctxErr == nil {
		return nil
	}

	if errors.Is(ctxErr, context.Canceled) {
		return ctxErr
	}

	if errors.Is(ctxErr, context.DeadlineExceeded) && lastErr != nil {
		return lastErr
	}

	return ctxErr
}

func RetryOn(ctx context.Context, call func(context.Context) error, opts ...RetryOption) error {
	options := defaultRetryConfig()
	for _, fn := range opts {
		fn(options)
	}

	ctx, cancel := context.WithTimeout(ctx, options.timeout)
	defer cancel()

	backoff := options.initBackoff
	var lastErr error

	for {
		if err := ctx.Err(); err != nil {
			return contextTerminationError(err, lastErr)
		}

		err := call(ctx)
		if err == nil {
			return nil
		}

		if !shouldRetry(err, options.retryOnFns) {
			return err
		}

		lastErr = err

		select {
		case <-ctx.Done():
			return contextTerminationError(ctx.Err(), lastErr)
		case <-time.After(backoff):
		}

		backoff *= 2
		if backoff > options.maxBackoff {
			backoff = options.maxBackoff
		}
	}
}
