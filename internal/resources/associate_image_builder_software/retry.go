// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

import "time"

const (
	retryTimeout     = 15 * time.Minute
	retryInitBackoff = 10 * time.Second
	retryMaxBackoff  = 2 * time.Minute

	waitAssociableRetryTimeout     = 30 * time.Minute
	waitAssociableRetryInitBackoff = 30 * time.Second
	waitAssociableRetryMaxBackoff  = 1 * time.Minute
)
