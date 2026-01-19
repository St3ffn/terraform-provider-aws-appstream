// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

import "time"

const (
	createRetryTimeout     = 15 * time.Minute
	createRetryInitBackoff = 10 * time.Second
	createRetryMaxBackoff  = 2 * time.Minute

	deleteRetryTimeout     = 45 * time.Minute
	deleteRetryInitBackoff = 30 * time.Second
	deleteRetryMaxBackoff  = 1 * time.Minute
)
