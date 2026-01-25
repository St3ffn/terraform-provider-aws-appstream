// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package testhelpers

import (
	"fmt"
	"strings"
)

func HCLStringList(values []string) string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = fmt.Sprintf("%q", v)
	}
	return strings.Join(out, ", ")
}
