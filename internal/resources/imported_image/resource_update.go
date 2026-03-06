// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package imported_image

import (
	"context"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

// Update reconciles mutable fields for imported images.
// Imported images are treated as immutable artifacts, so update only applies tags.
func (r *resource) Update(ctx context.Context, req tfresource.UpdateRequest, resp *tfresource.UpdateResponse) {
	var plan resourceModel
	var state resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if plan.ARN.IsNull() || plan.ARN.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update imported image because arn must be known.",
		)
		return
	}

	diff := newResourceDiff(state, plan)
	arn := plan.ARN.ValueString()

	if diff.RequiresTagApply() {
		_, tagDiags := r.tags.Apply(ctx, arn, plan.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	newState, diags := r.readImportedImage(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if ctx.Err() != nil {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
