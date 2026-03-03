// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package directory_config

import "time"

const (
	createRetryTimeout     = 2 * time.Minute
	createRetryInitBackoff = 2 * time.Second
	createRetryMaxBackoff  = 30 * time.Second

	updateRetryTimeout     = 2 * time.Minute
	updateRetryInitBackoff = 2 * time.Second
	updateRetryMaxBackoff  = 10 * time.Second
)
