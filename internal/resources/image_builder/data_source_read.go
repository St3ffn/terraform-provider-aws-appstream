// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (ds *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config dataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if config.Name.IsNull() || config.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Configuration",
			"Cannot read image builder because name must be set and known.",
		)
		return
	}

	name := config.Name.ValueString()

	out, err := ds.appstreamClient.DescribeImageBuilders(ctx, &awsappstream.DescribeImageBuildersInput{
		Names: []string{name},
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Image Builder Not Found",
				fmt.Sprintf("No image builder named %q was found.", name),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream Image Builder",
			fmt.Sprintf("Could not read image builder %q: %v", name, err),
		)
		return
	}

	if len(out.ImageBuilders) == 0 {
		resp.Diagnostics.AddError(
			"AWS AppStream Image Builder Not Found",
			fmt.Sprintf("No image builder named %q was found.", name),
		)
		return
	}

	imageBuilder := out.ImageBuilders[0]
	if imageBuilder.Name == nil {
		resp.Diagnostics.AddError(
			"Unexpected AWS Response",
			fmt.Sprintf("Image builder %q was returned without required identifiers.", name),
		)
		return
	}

	state := &dataSourceModel{
		ID:                          types.StringValue(aws.ToString(imageBuilder.Name)),
		Name:                        types.StringValue(aws.ToString(imageBuilder.Name)),
		ImageARN:                    util.StringOrNull(imageBuilder.ImageArn),
		InstanceType:                util.StringOrNull(imageBuilder.InstanceType),
		Description:                 util.StringOrNull(imageBuilder.Description),
		DisplayName:                 util.StringOrNull(imageBuilder.DisplayName),
		VPCConfig:                   flattenVPCConfig(ctx, imageBuilder.VpcConfig, &resp.Diagnostics),
		IAMRoleARN:                  util.StringOrNull(imageBuilder.IamRoleArn),
		EnableDefaultInternetAccess: util.BoolOrNull(imageBuilder.EnableDefaultInternetAccess),
		DomainJoinInfo:              flattenDomainJoinInfo(ctx, imageBuilder.DomainJoinInfo, &resp.Diagnostics),
		AppstreamAgentVersion:       util.StringOrNull(imageBuilder.AppstreamAgentVersion),
		AccessEndpoints:             flattenAccessEndpoints(ctx, imageBuilder.AccessEndpoints, &resp.Diagnostics),
		RootVolumeConfig:            flattenRootVolumeConfig(ctx, imageBuilder.RootVolumeConfig, &resp.Diagnostics),
		Tags:                        types.MapNull(types.StringType),
		ARN:                         util.StringOrNull(imageBuilder.Arn),
		CreatedTime:                 util.StringFromTime(imageBuilder.CreatedTime),
		Platform:                    types.StringValue(string(imageBuilder.Platform)),
		NetworkAccessConfiguration:  flattenNetworkAccessConfiguration(ctx, imageBuilder.NetworkAccessConfiguration, &resp.Diagnostics),
		LatestAppstreamAgentVersion: types.StringValue(string(imageBuilder.LatestAppstreamAgentVersion)),
		State:                       types.StringValue(string(imageBuilder.State)),
		StateChangeReason:           flattenStateChangeReason(ctx, imageBuilder.StateChangeReason, &resp.Diagnostics),
		ImageBuilderErrors:          flattenImageBuilderErrors(ctx, imageBuilder.ImageBuilderErrors, &resp.Diagnostics),
	}

	if !state.ARN.IsNull() {
		tags, diags := ds.tags.Read(ctx, state.ARN.ValueString())
		resp.Diagnostics.Append(diags...)
		state.Tags = tags
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
