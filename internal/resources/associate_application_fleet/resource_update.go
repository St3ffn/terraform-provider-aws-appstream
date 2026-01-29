// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_application_fleet

import (
	"context"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *resource) Update(_ context.Context, _ tfresource.UpdateRequest, _ *tfresource.UpdateResponse) {
	// no-op: all attributes require replacement
}
