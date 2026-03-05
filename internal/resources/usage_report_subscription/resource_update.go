// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription

import (
	"context"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

// Update is a no-op because this account-level singleton has no mutable arguments.
func (r *resource) Update(_ context.Context, _ tfresource.UpdateRequest, _ *tfresource.UpdateResponse) {
	// no-op: no attributes
}
