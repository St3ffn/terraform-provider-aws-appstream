// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet

import "time"

const (
	createRetryTimeout     = 15 * time.Minute
	createRetryInitBackoff = 10 * time.Second
	createRetryMaxBackoff  = 2 * time.Minute
)
