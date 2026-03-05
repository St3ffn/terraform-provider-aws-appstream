// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"errors"
)

// IsContextCanceled reports whether err was caused by context cancellation or deadline.
func IsContextCanceled(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
