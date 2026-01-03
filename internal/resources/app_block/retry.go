// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import "time"

const (
	createRetryTimeout     = 5 * time.Minute
	createRetryInitBackoff = 2 * time.Second
	createRetryMaxBackoff  = 30 * time.Second
)
