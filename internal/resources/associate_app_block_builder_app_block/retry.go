// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_app_block_builder_app_block

import "time"

const (
	createRetryTimeout     = 5 * time.Minute
	createRetryInitBackoff = 2 * time.Second
	createRetryMaxBackoff  = 1 * time.Minute
)
