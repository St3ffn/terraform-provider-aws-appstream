// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import (
	"context"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *resource) Update(ctx context.Context, req tfresource.UpdateRequest, resp *tfresource.UpdateResponse) {
	var plan model
	var state model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if plan.ID.IsNull() || plan.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update app block because id must be known.",
		)
		return
	}

	arn := plan.ID.ValueString()

	// guard against unexpected identity drift
	if !state.ID.IsNull() && !state.ID.IsUnknown() {
		if state.ID.ValueString() != arn {
			resp.Diagnostics.AddError(
				"Unexpected Update Request",
				"App block identity (ARN) changed during update. This should trigger replacement. Please report this issue.",
			)
			return
		}
	}

	_, tagDiags := r.tags.Apply(ctx, arn, plan.Tags)
	resp.Diagnostics.Append(tagDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readAppBlock(ctx, plan)
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
