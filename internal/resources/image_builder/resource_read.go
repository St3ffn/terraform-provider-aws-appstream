// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

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
	var state resourceModel

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

	newState, diags := r.readImageBuilder(ctx, state)
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

func (r *resource) readImageBuilder(ctx context.Context, prior resourceModel) (*resourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	name := prior.Name.ValueString()

	out, err := r.appstreamClient.DescribeImageBuilders(ctx, &awsappstream.DescribeImageBuildersInput{
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
			"Error Reading AWS AppStream Image Builder",
			fmt.Sprintf("Could not read image builder %q: %v", name, err),
		)
		return nil, diags
	}

	if len(out.ImageBuilders) == 0 {
		return nil, diags
	}

	imageBuilder := out.ImageBuilders[0]
	if imageBuilder.Name == nil {
		return nil, diags
	}

	state := &resourceModel{
		ID:                          types.StringValue(aws.ToString(imageBuilder.Name)),
		Name:                        types.StringValue(aws.ToString(imageBuilder.Name)),
		ImageName:                   prior.ImageName,
		ImageARN:                    util.StringOrNull(imageBuilder.ImageArn),
		InstanceType:                util.StringOrNull(imageBuilder.InstanceType),
		Description:                 util.StringOrNull(imageBuilder.Description),
		DisplayName:                 util.StringOrNull(imageBuilder.DisplayName),
		VPCConfig:                   flattenVPCConfig(ctx, imageBuilder.VpcConfig, &diags),
		IAMRoleARN:                  util.StringOrNull(imageBuilder.IamRoleArn),
		EnableDefaultInternetAccess: util.BoolOrNull(imageBuilder.EnableDefaultInternetAccess),
		DomainJoinInfo:              flattenDomainJoinInfo(ctx, imageBuilder.DomainJoinInfo, &diags),
		AppstreamAgentVersion:       util.StringOrNull(imageBuilder.AppstreamAgentVersion),
		AccessEndpoints:             flattenAccessEndpoints(ctx, imageBuilder.AccessEndpoints, &diags),
		RootVolumeConfig:            flattenRootVolumeConfig(ctx, imageBuilder.RootVolumeConfig, &diags),
		Tags:                        types.MapNull(types.StringType),
		ARN:                         util.StringOrNull(imageBuilder.Arn),
		CreatedTime:                 util.StringFromTime(imageBuilder.CreatedTime),
		Platform:                    types.StringValue(string(imageBuilder.Platform)),
		NetworkAccessConfiguration:  flattenNetworkAccessConfiguration(ctx, imageBuilder.NetworkAccessConfiguration, &diags),
		LatestAppstreamAgentVersion: types.StringValue(string(imageBuilder.LatestAppstreamAgentVersion)),
		State:                       types.StringValue(string(imageBuilder.State)),
		StateChangeReason:           flattenStateChangeReason(ctx, imageBuilder.StateChangeReason, &diags),
		ImageBuilderErrors:          flattenImageBuilderErrors(ctx, imageBuilder.ImageBuilderErrors, &diags),
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
