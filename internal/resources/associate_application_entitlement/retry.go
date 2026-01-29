// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_application_entitlement

import "time"

const (
	createRetryTimeout     = 3 * time.Minute
	createRetryInitBackoff = 2 * time.Second
	createRetryMaxBackoff  = 30 * time.Second
)
