// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_application_entitlement

import (
	"context"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

// Update is a no-op because all association attributes are ForceNew.
func (r *resource) Update(_ context.Context, _ tfresource.UpdateRequest, _ *tfresource.UpdateResponse) {
	// no-op: all attributes require replacement
}
