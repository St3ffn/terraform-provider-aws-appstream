// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

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
	var config model

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
			"Cannot read app block builder because name must be set and known.",
		)
		return
	}

	name := config.Name.ValueString()

	out, err := ds.appstreamClient.DescribeAppBlockBuilders(ctx, &awsappstream.DescribeAppBlockBuildersInput{
		Names: []string{name},
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream App Block Builder Not Found",
				fmt.Sprintf("No app block builder named %q was found.", name),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream App Block Builder",
			fmt.Sprintf("Could not read app block builder %q: %v", name, err),
		)
		return
	}

	if len(out.AppBlockBuilders) == 0 {
		resp.Diagnostics.AddError(
			"AWS AppStream App Block Builder Not Found",
			fmt.Sprintf("No app block builder named %q was found.", name),
		)
		return
	}

	appBlockBuilder := out.AppBlockBuilders[0]
	if appBlockBuilder.Name == nil {
		resp.Diagnostics.AddError(
			"Unexpected AWS Response",
			fmt.Sprintf("App block builder %q was returned without required identifiers.", name),
		)
		return
	}

	state := &model{
		ID:                          types.StringValue(aws.ToString(appBlockBuilder.Name)),
		Name:                        types.StringValue(aws.ToString(appBlockBuilder.Name)),
		InstanceType:                util.StringOrNull(appBlockBuilder.InstanceType),
		Platform:                    types.StringValue(string(appBlockBuilder.Platform)),
		Description:                 util.StringOrNull(appBlockBuilder.Description),
		DisplayName:                 util.StringOrNull(appBlockBuilder.DisplayName),
		VPCConfig:                   flattenVPCConfig(ctx, appBlockBuilder.VpcConfig, &resp.Diagnostics),
		IAMRoleARN:                  util.StringOrNull(appBlockBuilder.IamRoleArn),
		EnableDefaultInternetAccess: util.BoolOrNull(appBlockBuilder.EnableDefaultInternetAccess),
		AccessEndpoints:             flattenAccessEndpoints(ctx, appBlockBuilder.AccessEndpoints, &resp.Diagnostics),
		Tags:                        types.MapNull(types.StringType),
		ARN:                         util.StringOrNull(appBlockBuilder.Arn),
		CreatedTime:                 util.StringFromTime(appBlockBuilder.CreatedTime),
		State:                       types.StringValue(string(appBlockBuilder.State)),
		AppBlockBuilderErrors:       flattenAppBlockBuilderErrors(ctx, appBlockBuilder.AppBlockBuilderErrors, &resp.Diagnostics),
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
