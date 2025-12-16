// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *entitlementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan entitlementModel
	var state entitlementModel

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

	var description *string
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		description = plan.Description.ValueStringPointer()
	}

	_, err := r.appStreamClient.UpdateEntitlement(ctx, &awsappstream.UpdateEntitlementInput{
		StackName:     aws.String(stackName),
		Name:          aws.String(name),
		Description:   description,
		AppVisibility: awstypes.AppVisibility(plan.AppVisibility.ValueString()),
		Attributes:    awsAttrs,
	})
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}

		// disappeared, treat as gone and next plan/apply will recreate
		if isAppStreamNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Updating AWS AppStream Entitlement",
			fmt.Sprintf("Could not update entitlement %q in stack %q: %v", name, stackName, err),
		)
		return
	}

	plan.ID = types.StringValue(buildEntitlementID(stackName, name))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
