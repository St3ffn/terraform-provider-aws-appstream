// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package imported_image

import "time"

const (
	createRetryTimeout     = 10 * time.Minute
	createRetryInitBackoff = 5 * time.Second
	createRetryMaxBackoff  = 30 * time.Second

	waitAvailableRetryTimeout     = 90 * time.Minute
	waitAvailableRetryInitBackoff = 30 * time.Second
	waitAvailableRetryMaxBackoff  = 1 * time.Minute
)
