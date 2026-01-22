// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

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

	if state.Name.IsNull() || state.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Required attribute name is missing from state. "+
				"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
		)
		return
	}

	newState, diags := r.readAppBlockBuilder(ctx, state)
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

func (r *resource) readAppBlockBuilder(ctx context.Context, prior model) (*model, diag.Diagnostics) {
	var diags diag.Diagnostics

	name := prior.Name.ValueString()

	out, err := r.appstreamClient.DescribeAppBlockBuilders(ctx, &awsappstream.DescribeAppBlockBuildersInput{
		Names: []string{name},
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		if util.IsAppStreamNotFound(err) {
			return nil, diags
		}

		diags.AddError(
			"Error Reading AWS AppStream App Block Builder",
			fmt.Sprintf("Could not read app block builder %q: %v", name, err),
		)
		return nil, diags
	}

	if len(out.AppBlockBuilders) == 0 {
		return nil, diags
	}

	appBlockBuilder := out.AppBlockBuilders[0]
	if appBlockBuilder.Name == nil {
		return nil, diags
	}

	state := &model{
		ID:                          types.StringValue(aws.ToString(appBlockBuilder.Name)),
		Name:                        types.StringValue(aws.ToString(appBlockBuilder.Name)),
		InstanceType:                util.StringOrNull(appBlockBuilder.InstanceType),
		Platform:                    types.StringValue(string(appBlockBuilder.Platform)),
		Description:                 util.StringOrNull(appBlockBuilder.Description),
		DisplayName:                 util.StringOrNull(appBlockBuilder.DisplayName),
		VPCConfig:                   flattenVPCConfig(ctx, appBlockBuilder.VpcConfig, &diags),
		IAMRoleARN:                  util.StringOrNull(appBlockBuilder.IamRoleArn),
		EnableDefaultInternetAccess: util.BoolOrNull(appBlockBuilder.EnableDefaultInternetAccess),
		AccessEndpoints:             flattenAccessEndpoints(ctx, appBlockBuilder.AccessEndpoints, &diags),
		Tags:                        types.MapNull(types.StringType),
		ARN:                         util.StringOrNull(appBlockBuilder.Arn),
		CreatedTime:                 util.StringFromTime(appBlockBuilder.CreatedTime),
		State:                       types.StringValue(string(appBlockBuilder.State)),
		StateChangeReason:           flattenStateChangeReason(ctx, appBlockBuilder.StateChangeReason, &diags),
		AppBlockBuilderErrors:       flattenAppBlockBuilderErrors(ctx, appBlockBuilder.AppBlockBuilderErrors, &diags),
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
