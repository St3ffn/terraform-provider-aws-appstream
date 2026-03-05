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

// RetryOnFn classifies whether an error should be retried.
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

// RetryOption mutates retry configuration used by RetryOn and RetryOnValue.
type RetryOption func(*retryConfig)

// WithTimeout sets the maximum duration for the retry operation.
func WithTimeout(timeout time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.timeout = timeout
	}
}

// WithInitBackoff sets the initial retry backoff duration.
func WithInitBackoff(initBackoff time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.initBackoff = initBackoff
	}
}

// WithMaxBackoff sets the maximum retry backoff duration.
func WithMaxBackoff(maxBackoff time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.maxBackoff = maxBackoff
	}
}

// WithRetryOnFns sets the retry predicate functions used to classify retryable errors.
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

// RetryOn executes call and retries only when the returned error matches any
// predicate from WithRetryOnFns. Retry backoff starts at init_backoff, doubles per
// attempt, and is capped by max_backoff until timeout/context cancellation.
//
// Return behavior:
//   - success: nil
//   - non-retryable error: immediate return of that error
//   - canceled context: context.Canceled
//   - timeout: last retryable error if available, otherwise context.DeadlineExceeded
func RetryOn(ctx context.Context, call func(context.Context) error, opts ...RetryOption) error {
	_, err := RetryOnValue(
		ctx,
		func(ctx context.Context) (struct{}, error) {
			return struct{}{}, call(ctx)
		},
		opts...,
	)
	return err
}

// RetryOnValue is RetryOn with a typed return value from call.
// The last successful output is returned together with nil on success, and the
// zero value of T is returned on failures before call succeeds.
func RetryOnValue[T any](ctx context.Context, call func(context.Context) (T, error), opts ...RetryOption) (T, error) {
	options := defaultRetryConfig()
	for _, fn := range opts {
		fn(options)
	}

	ctx, cancel := context.WithTimeout(ctx, options.timeout)
	defer cancel()

	backoff := options.initBackoff
	var (
		out     T
		lastErr error
	)

	for {
		if err := ctx.Err(); err != nil {
			return out, contextTerminationError(err, lastErr)
		}

		out, err := call(ctx)
		if err == nil {
			return out, nil
		}

		if !shouldRetry(err, options.retryOnFns) {
			return out, err
		}

		lastErr = err

		select {
		case <-ctx.Done():
			return out, contextTerminationError(ctx.Err(), lastErr)
		case <-time.After(backoff):
		}

		backoff *= 2
		if backoff > options.maxBackoff {
			backoff = options.maxBackoff
		}
	}
}
