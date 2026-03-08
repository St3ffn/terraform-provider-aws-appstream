// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package imported_image

import "time"

const (
	createRetryTimeout     = 5 * time.Minute
	createRetryInitBackoff = 5 * time.Second
	createRetryMaxBackoff  = 30 * time.Second

	waitAvailableRetryTimeout     = 3 * time.Hour
	waitAvailableRetryInitBackoff = 30 * time.Second
	waitAvailableRetryMaxBackoff  = 1 * time.Minute
)
