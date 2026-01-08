// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

import "time"

const (
	createRetryTimeout     = 15 * time.Minute
	createRetryInitBackoff = 10 * time.Second
	createRetryMaxBackoff  = 2 * time.Minute

	imageBuilderWaitTimeout     = 45 * time.Minute
	imageBuilderWaitInitBackoff = 30 * time.Second
	imageBuilderWaitMaxBackoff  = 1 * time.Minute
)
