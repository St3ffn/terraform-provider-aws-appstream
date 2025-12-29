// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

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

	if plan.ID.IsNull() || plan.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update application because id must be known.",
		)
		return
	}

	arn := plan.ID.ValueString()

	name, err := applicationNameFromARN(arn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			fmt.Sprintf("Could not parse application name from ARN %q: %v", arn, err),
		)
		return
	}

	// guard against unexpected identity drift
	if !state.ID.IsNull() && !state.ID.IsUnknown() {
		if state.ID.ValueString() != arn {
			resp.Diagnostics.AddError(
				"Unexpected Update Request",
				"Application identity (ARN) changed during update. This should trigger replacement. Please report this issue.",
			)
			return
		}
	}

	input := &awsappstream.UpdateApplicationInput{
		Name: aws.String(name),
	}

	var attrsToDelete []awstypes.ApplicationAttribute

	util.OptionalStringUpdate(plan.DisplayName, state.DisplayName, func(v *string) {
		input.DisplayName = v
	})
	util.OptionalStringUpdate(plan.Description, state.Description, func(v *string) {
		input.Description = v
	})
	util.OptionalStringUpdate(plan.LaunchPath, state.LaunchPath, func(v *string) {
		input.LaunchPath = v
	})

	util.OptionalStringUpdate(plan.AppBlockARN, state.AppBlockARN, func(v *string) {
		input.AppBlockArn = v
	})

	if !plan.WorkingDirectory.IsUnknown() {
		if plan.WorkingDirectory.IsNull() {
			attrsToDelete = append(attrsToDelete, awstypes.ApplicationAttributeWorkingDirectory)
		} else {
			input.WorkingDirectory = plan.WorkingDirectory.ValueStringPointer()
		}
	}

	if !plan.LaunchParameters.IsUnknown() {
		if plan.LaunchParameters.IsNull() {
			attrsToDelete = append(attrsToDelete, awstypes.ApplicationAttributeLaunchParameters)
		} else {
			input.LaunchParameters = plan.LaunchParameters.ValueStringPointer()
		}
	}

	if !plan.IconS3Location.IsUnknown() {
		if plan.IconS3Location.IsNull() {
			// no delete support
		} else {
			input.IconS3Location = expandIconS3Location(
				ctx, plan.IconS3Location, &resp.Diagnostics,
			)
		}
	}

	if len(attrsToDelete) > 0 {
		input.AttributesToDelete = attrsToDelete
	}

	if resp.Diagnostics.HasError() {
		return
	}

	_, err = r.appstreamClient.UpdateApplication(ctx, input)
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Updating AWS AppStream Application",
			fmt.Sprintf("Could not update application %q: %v", arn, err),
		)
		return
	}

	_, tagDiags := r.tags.Apply(ctx, arn, plan.Tags)
	resp.Diagnostics.Append(tagDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readApplication(ctx, arn)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
