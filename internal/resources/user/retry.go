// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user

import "time"

const (
	createRetryTimeout     = 5 * time.Minute
	createRetryInitBackoff = 2 * time.Second
	createRetryMaxBackoff  = 30 * time.Second

	disableRetryTimeout     = 3 * time.Minute
	disableRetryInitBackoff = 2 * time.Second
	disableRetryMaxBackoff  = 30 * time.Second

	readRetryTimeout     = 2 * time.Minute
	readRetryInitBackoff = 2 * time.Second
	readRetryMaxBackoff  = 20 * time.Second
)
