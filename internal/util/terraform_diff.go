// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

func Changed(state, plan attr.Value) bool {
	if plan.IsUnknown() {
		return false
	}
	return !state.Equal(plan)
}
