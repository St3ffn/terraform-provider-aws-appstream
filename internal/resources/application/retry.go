// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

import "time"

const (
	createRetryTimeout     = 10 * time.Minute
	createRetryInitBackoff = 5 * time.Second
	createRetryMaxBackoff  = 1 * time.Minute
)
