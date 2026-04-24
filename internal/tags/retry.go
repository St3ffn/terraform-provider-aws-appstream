// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package tags

import "time"

const (
	taggingRetryTimeout     = 5 * time.Minute
	taggingRetryInitBackoff = 2 * time.Second
	taggingRetryMaxBackoff  = 30 * time.Second
)
