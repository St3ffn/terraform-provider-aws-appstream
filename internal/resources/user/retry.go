// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package user

import "time"

const (
	disableRetryTimeout     = 3 * time.Minute
	disableRetryInitBackoff = 2 * time.Second
	disableRetryMaxBackoff  = 30 * time.Second

	readRetryTimeout     = 2 * time.Minute
	readRetryInitBackoff = 2 * time.Second
	readRetryMaxBackoff  = 20 * time.Second
)
