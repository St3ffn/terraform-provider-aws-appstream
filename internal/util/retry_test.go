// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDefaultRetryConfig(t *testing.T) {
	cfg := defaultRetryConfig()

	if cfg.timeout != defaultRetryTimeout {
		t.Errorf("timeout = %v, want %v", cfg.timeout, defaultRetryTimeout)
	}

	if cfg.initBackoff != defaultRetryInitBackoff {
		t.Errorf("initBackoff = %v, want %v", cfg.initBackoff, defaultRetryInitBackoff)
	}

	if cfg.maxBackoff != defaultRetryMaxBackoff {
		t.Errorf("maxBackoff = %v, want %v", cfg.maxBackoff, defaultRetryMaxBackoff)
	}

	if len(cfg.retryOnFns) != 0 {
		t.Errorf("retryOnFns must be empty")
	}
}

func TestRetryOn_SuccessFirstTry(t *testing.T) {
	ctx := context.Background()

	calls := 0
	err := RetryOn(ctx, func(context.Context) error {
		calls++
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if calls != 1 {
		t.Fatalf("call count = %d, want 1", calls)
	}
}

func TestRetryOn_RetryThenSuccess(t *testing.T) {
	ctx := context.Background()

	calls := 0
	retryErr := errors.New("retry")

	err := RetryOn(
		ctx,
		func(context.Context) error {
			calls++
			if calls < 3 {
				return retryErr
			}
			return nil
		},
		WithRetryOnFns(func(err error) bool {
			return errors.Is(err, retryErr)
		}),
		WithInitBackoff(1*time.Millisecond),
		WithMaxBackoff(2*time.Millisecond),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if calls != 3 {
		t.Fatalf("call count = %d, want 3", calls)
	}
}

func TestRetryOn_NonRetryableError(t *testing.T) {
	ctx := context.Background()

	calls := 0
	expectedErr := errors.New("fail fast")

	err := RetryOn(
		ctx,
		func(context.Context) error {
			calls++
			return expectedErr
		},
		WithRetryOnFns(func(error) bool { return false }),
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("error = %v, want %v", err, expectedErr)
	}

	if calls != 1 {
		t.Fatalf("call count = %d, want 1", calls)
	}
}

func TestRetryOn_ContextTimeout(t *testing.T) {
	ctx := context.Background()

	lastErr := errors.New("last_error")
	err := RetryOn(
		ctx,
		func(context.Context) error {
			return lastErr
		},
		WithRetryOnFns(func(error) bool { return true }),
		WithTimeout(5*time.Millisecond),
		WithInitBackoff(1*time.Millisecond),
	)

	if !errors.Is(err, lastErr) {
		t.Fatalf("error = %v, want last_error", err)
	}
}
