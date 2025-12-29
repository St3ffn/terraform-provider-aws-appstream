// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package testhelpers

import (
	"os"
	"testing"
)

func TestAccPreCheck(t *testing.T) {
	t.Helper()

	if os.Getenv("AWS_PROFILE") == "" &&
		os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Fatal("AWS credentials not set")
	}

	if os.Getenv("AWS_REGION") == "" &&
		os.Getenv("AWS_DEFAULT_REGION") == "" {
		t.Fatal("AWS region not set")
	}
}
