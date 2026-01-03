// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack

import (
	"context"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *resource) Update(_ context.Context, _ tfresource.UpdateRequest, _ *tfresource.UpdateResponse) {
	// no-op: all attributes require replacement
}
