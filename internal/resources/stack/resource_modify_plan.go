// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"context"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *resource) ModifyPlan(ctx context.Context, req tfresource.ModifyPlanRequest, resp *tfresource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.TagsAll = r.tags.EffectiveTagsForPlan(plan.Tags)

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}
