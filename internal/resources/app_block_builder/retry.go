// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import "time"

const (
	createRetryTimeout     = 15 * time.Minute
	createRetryInitBackoff = 10 * time.Second
	createRetryMaxBackoff  = 1 * time.Minute

	deleteRetryTimeout     = 45 * time.Minute
	deleteRetryInitBackoff = 30 * time.Second
	deleteRetryMaxBackoff  = 1 * time.Minute

	updateRetryTimeout     = 45 * time.Minute
	updateRetryInitBackoff = 30 * time.Second
	updateRetryMaxBackoff  = 1 * time.Minute
)
