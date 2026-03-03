// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_theme

import "time"

const (
	createRetryTimeout     = 10 * time.Minute
	createRetryInitBackoff = 5 * time.Second
	createRetryMaxBackoff  = 30 * time.Second

	updateRetryTimeout     = 5 * time.Minute
	updateRetryInitBackoff = 2 * time.Second
	updateRetryMaxBackoff  = 30 * time.Second
)
