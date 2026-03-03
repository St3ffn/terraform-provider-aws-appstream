// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_fleet_stack

import "time"

const (
	createRetryTimeout     = 5 * time.Minute
	createRetryInitBackoff = 2 * time.Second
	createRetryMaxBackoff  = 30 * time.Second

	deleteRetryTimeout     = 5 * time.Minute
	deleteRetryInitBackoff = 2 * time.Second
	deleteRetryMaxBackoff  = 30 * time.Second
)
