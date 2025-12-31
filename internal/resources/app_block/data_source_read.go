// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

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

	if config.ARN.IsNull() || config.ARN.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Configuration",
			"Cannot read app block because arn must be set and known.",
		)
		return
	}

	arn := config.ARN.ValueString()

	out, err := ds.appstreamClient.DescribeAppBlocks(ctx, &awsappstream.DescribeAppBlocksInput{
		Arns: []string{arn},
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream App Block Not Found",
				fmt.Sprintf("No app block with ARN %q was found.", arn),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream App Block",
			fmt.Sprintf("Could not read app block %q: %v", arn, err),
		)
		return
	}

	if len(out.AppBlocks) == 0 {
		resp.Diagnostics.AddError(
			"AWS AppStream App Block Not Found",
			fmt.Sprintf("No app block with ARN %q was found.", arn),
		)
		return
	}

	appBlock := out.AppBlocks[0]
	if appBlock.Arn == nil || appBlock.Name == nil {
		resp.Diagnostics.AddError(
			"Unexpected AWS Response",
			fmt.Sprintf("App block %q was returned without required identifiers.", arn),
		)
		return
	}

	state := &model{
		ID:                     types.StringValue(aws.ToString(appBlock.Arn)),
		ARN:                    types.StringValue(aws.ToString(appBlock.Arn)),
		Name:                   types.StringValue(aws.ToString(appBlock.Name)),
		DisplayName:            util.StringOrNull(appBlock.DisplayName),
		Description:            util.StringOrNull(appBlock.Description),
		SourceS3Location:       flattenSourceS3LocationData(ctx, appBlock.SourceS3Location, &resp.Diagnostics),
		SetupScriptDetails:     flattenScriptDetailsData(ctx, appBlock.SetupScriptDetails, &resp.Diagnostics),
		PostSetupScriptDetails: flattenScriptDetailsData(ctx, appBlock.PostSetupScriptDetails, &resp.Diagnostics),
		PackagingType:          types.StringValue(string(appBlock.PackagingType)),
		Tags:                   types.Map{},
		CreatedTime:            util.StringFromTime(appBlock.CreatedTime),
		AppBlockErrors:         flattenAppBlockErrorsData(ctx, appBlock.AppBlockErrors, &resp.Diagnostics),
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
