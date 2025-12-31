// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	// Validate required attributes
	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create app block because name must be known.",
		)
		return
	}

	name := plan.Name.ValueString()

	input := &awsappstream.CreateAppBlockInput{
		Name: aws.String(name),
	}

	input.DisplayName = util.StringPointerOrNil(plan.DisplayName)
	input.Description = util.StringPointerOrNil(plan.Description)

	if !plan.PackagingType.IsNull() && !plan.PackagingType.IsUnknown() {
		input.PackagingType = awstypes.PackagingType(plan.PackagingType.ValueString())
	}

	if !plan.SourceS3Location.IsNull() && !plan.SourceS3Location.IsUnknown() {
		input.SourceS3Location = expandSourceS3Location(ctx, plan.SourceS3Location, &resp.Diagnostics)
	}

	if !plan.SetupScriptDetails.IsNull() && !plan.SetupScriptDetails.IsUnknown() {
		input.SetupScriptDetails = expandScriptDetails(ctx, plan.SetupScriptDetails, &resp.Diagnostics)
	}

	if !plan.PostSetupScriptDetails.IsNull() && !plan.PostSetupScriptDetails.IsUnknown() {
		input.PostSetupScriptDetails = expandScriptDetails(ctx, plan.PostSetupScriptDetails, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.appstreamClient.CreateAppBlock(ctx, input)
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsResourceAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream App Block Already Exists",
				fmt.Sprintf(
					"An app block named %q already exists. "+
						"To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> <app_block_arn>",
					name,
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream App Block",
			fmt.Sprintf("Could not create app block %q: %v", name, err),
		)
		return
	}

	if out != nil && out.AppBlock != nil && out.AppBlock.Arn != nil {
		_, tagDiags := r.tags.Apply(ctx, aws.ToString(out.AppBlock.Arn), plan.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	newState, diags := r.readAppBlock(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if util.IsContextCanceled(ctx.Err()) {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
