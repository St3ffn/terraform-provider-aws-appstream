// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package entitlement

import "time"

const (
	createRetryTimeout     = 5 * time.Minute
	createRetryInitBackoff = 2 * time.Second
	createRetryMaxBackoff  = 1 * time.Minute
)
