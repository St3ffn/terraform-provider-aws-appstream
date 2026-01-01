// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
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

	if plan.StackName.IsNull() || plan.StackName.IsUnknown() ||
		plan.Name.IsNull() || plan.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update entitlement because stack_name and name must be known.",
		)
		return
	}

	stackName := plan.StackName.ValueString()
	name := plan.Name.ValueString()

	// guard against unexpected identity drift
	if !state.StackName.IsNull() && !state.StackName.IsUnknown() &&
		!state.Name.IsNull() && !state.Name.IsUnknown() {
		if state.StackName.ValueString() != stackName || state.Name.ValueString() != name {
			resp.Diagnostics.AddError(
				"Unexpected Update Request",
				"Entitlement identity (stack_name/name) changed during update. This should trigger replacement. Please report this issue.",
			)
			return
		}
	}

	if plan.AppVisibility.IsNull() || plan.AppVisibility.IsUnknown() ||
		plan.Attributes.IsNull() || plan.Attributes.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update entitlement because app_visibility and attributes must be known.",
		)
		return
	}

	awsAttrs := expandEntitlementAttributes(ctx, plan.Attributes, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &awsappstream.UpdateEntitlementInput{
		StackName:     aws.String(stackName),
		Name:          aws.String(name),
		AppVisibility: awstypes.AppVisibility(plan.AppVisibility.ValueString()),
		Attributes:    awsAttrs,
	}

	util.OptionalStringUpdate(plan.Description, state.Description, func(description *string) {
		input.Description = description
	})

	_, err := r.appstreamClient.UpdateEntitlement(ctx, input)
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		// disappeared, treat as gone and next plan/apply will recreate
		if util.IsAppStreamNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Updating AWS AppStream Entitlement",
			fmt.Sprintf("Could not update entitlement %q in stack %q: %v", name, stackName, err),
		)
		return
	}

	newState, diags := r.readEntitlement(ctx, plan)
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
