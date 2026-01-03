// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

import "time"

const (
	createRetryTimeout     = 2 * time.Minute
	createRetryInitBackoff = 2 * time.Second
	createRetryMaxBackoff  = 20 * time.Second
)
