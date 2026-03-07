// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package copied_image

import (
	"context"

	awstaggingapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

// Update reconciles mutable fields for copied images.
// Copied images are treated as immutable artifacts, so update only applies tags.
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

	if plan.ARN.IsNull() || plan.ARN.IsUnknown() ||
		plan.DestinationRegion.IsNull() || plan.DestinationRegion.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update copied image because arn and destination_region must be known.",
		)
		return
	}

	diff := newResourceDiff(state, plan)
	arn := plan.ARN.ValueString()
	destinationRegion := plan.DestinationRegion.ValueString()

	if diff.RequiresTagApply() {
		_, tagDiags := r.tags.Apply(
			ctx,
			arn,
			plan.Tags,
			func(o *awstaggingapi.Options) {
				o.Region = destinationRegion
			},
		)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	newState, diags := r.readCopiedImage(ctx, plan)
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
