// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder

import "time"

const (
	createRetryTimeout     = 15 * time.Minute
	createRetryInitBackoff = 10 * time.Second
	createRetryMaxBackoff  = 1 * time.Minute

	deleteRetryTimeout     = 45 * time.Minute
	deleteRetryInitBackoff = 30 * time.Second
	deleteRetryMaxBackoff  = 1 * time.Minute
)
