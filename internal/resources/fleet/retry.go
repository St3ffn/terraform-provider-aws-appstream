// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package fleet

import "time"

const (
	createRetryTimeout     = 15 * time.Minute
	createRetryInitBackoff = 10 * time.Second
	createRetryMaxBackoff  = 1 * time.Minute

	updateRetryTimeout     = 45 * time.Minute
	updateRetryInitBackoff = 30 * time.Second
	updateRetryMaxBackoff  = 1 * time.Minute
)
