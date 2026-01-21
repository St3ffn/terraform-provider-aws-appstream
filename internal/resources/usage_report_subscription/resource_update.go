// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription

import (
	"context"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *resource) Update(_ context.Context, _ tfresource.UpdateRequest, _ *tfresource.UpdateResponse) {
	// no-op: no attributes
}
