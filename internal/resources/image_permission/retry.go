// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_permission

import "time"

const (
	retryTimeout     = 5 * time.Minute
	retryInitBackoff = 5 * time.Second
	retryMaxBackoff  = 30 * time.Second
)
