// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import "time"

const (
	createRetryTimeout     = 15 * time.Minute
	createRetryInitBackoff = 10 * time.Second
	createRetryMaxBackoff  = 1 * time.Minute

	deleteRetryTimeout     = 45 * time.Minute
	deleteRetryInitBackoff = 30 * time.Second
	deleteRetryMaxBackoff  = 1 * time.Minute

	stopRetryTimeout     = 45 * time.Minute
	stopRetryInitBackoff = 30 * time.Second
	stopRetryMaxBackoff  = 1 * time.Minute

	startRetryTimeout     = 45 * time.Minute
	startRetryInitBackoff = 30 * time.Second
	startRetryMaxBackoff  = 1 * time.Minute

	updateRetryTimeout     = 5 * time.Minute
	updateRetryInitBackoff = 2 * time.Second
	updateRetryMaxBackoff  = 30 * time.Second
)
