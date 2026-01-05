// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Read(ctx context.Context, req tfresource.ReadRequest, resp *tfresource.ReadResponse) {
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if state.ID.IsNull() || state.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Required attribute id is missing from state. "+
				"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
		)
		return
	}

	newState, diags := r.readAppBlock(ctx, state)
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

func (r *resource) readAppBlock(ctx context.Context, prior model) (*model, diag.Diagnostics) {
	var diags diag.Diagnostics

	arn := prior.ID.ValueString()

	out, err := r.appstreamClient.DescribeAppBlocks(ctx, &awsappstream.DescribeAppBlocksInput{
		Arns: []string{arn},
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		if util.IsAppStreamNotFound(err) {
			return nil, diags
		}

		diags.AddError(
			"Error Reading AWS AppStream App Block",
			fmt.Sprintf("Could not read app block %q: %v", arn, err),
		)
		return nil, diags
	}

	if len(out.AppBlocks) == 0 {
		return nil, diags
	}

	appBlock := out.AppBlocks[0]
	if appBlock.Arn == nil || appBlock.Name == nil {
		return nil, diags
	}

	state := &model{
		ID:                     types.StringValue(aws.ToString(appBlock.Arn)),
		Name:                   types.StringValue(aws.ToString(appBlock.Name)),
		DisplayName:            util.StringOrNull(appBlock.DisplayName),
		Description:            util.StringOrNull(appBlock.Description),
		SourceS3Location:       flattenSourceS3LocationData(ctx, appBlock.SourceS3Location, &diags),
		SetupScriptDetails:     flattenScriptDetailsResource(ctx, prior.SetupScriptDetails, appBlock.SetupScriptDetails, &diags),
		PostSetupScriptDetails: flattenScriptDetailsResource(ctx, prior.PostSetupScriptDetails, appBlock.PostSetupScriptDetails, &diags),
		PackagingType:          util.FlattenStateOwnedString(prior.PackagingType, aws.String(string(appBlock.PackagingType))),
		Tags:                   types.MapNull(types.StringType),
		ARN:                    util.StringOrNull(appBlock.Arn),
		CreatedTime:            util.StringFromTime(appBlock.CreatedTime),
		AppBlockErrors:         flattenAppBlockErrorsData(ctx, appBlock.AppBlockErrors, &diags),
	}

	if !state.ARN.IsNull() {
		tags, tagDiags := r.tags.Read(ctx, state.ARN.ValueString())
		diags.Append(tagDiags...)
		state.Tags = tags
	}

	if diags.HasError() {
		return nil, diags
	}

	return state, diags
}
