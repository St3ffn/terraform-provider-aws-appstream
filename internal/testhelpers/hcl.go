// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package testhelpers

import (
	"fmt"
	"strings"
)

// HCLStringList formats a slice of strings as a comma-separated, quoted list
// suitable for inline use in Terraform HCL.
//
// Example:
//
//	subnets := []string{"subnet-123", "subnet-456"}
//	HCLStringList(subnets)
//	// => "subnet-123", "subnet-456"
//
// This helper is primarily used in acceptance tests when embedding dynamic
// string lists (such as subnet IDs) directly into Terraform configuration
// templates.
func HCLStringList(values []string) string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = fmt.Sprintf("%q", v)
	}
	return strings.Join(out, ", ")
}
