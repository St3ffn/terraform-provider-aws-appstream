// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

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
			"Cannot read application because arn must be set and known.",
		)
		return
	}

	arn := config.ARN.ValueString()

	out, err := ds.appstreamClient.DescribeApplications(ctx, &awsappstream.DescribeApplicationsInput{
		Arns: []string{arn},
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Application Not Found",
				fmt.Sprintf("No application with ARN %q was found.", arn),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream Application",
			fmt.Sprintf("Could not read application %q: %v", arn, err),
		)
		return
	}

	if len(out.Applications) == 0 {
		resp.Diagnostics.AddError(
			"AWS AppStream Application Not Found",
			fmt.Sprintf("No application with ARN %q was found.", arn),
		)
		return
	}

	app := out.Applications[0]
	if app.Arn == nil || app.Name == nil {
		resp.Diagnostics.AddError(
			"Unexpected AWS Response",
			fmt.Sprintf("Application %q was returned without required identifiers.", arn),
		)
		return
	}

	state := &model{
		ID:               types.StringValue(aws.ToString(app.Arn)),
		ARN:              types.StringValue(aws.ToString(app.Arn)),
		Name:             types.StringValue(aws.ToString(app.Name)),
		DisplayName:      util.StringOrNull(app.DisplayName),
		Description:      util.StringOrNull(app.Description),
		IconS3Location:   flattenIconS3Location(ctx, app.IconS3Location, &resp.Diagnostics),
		LaunchPath:       util.StringOrNull(app.LaunchPath),
		WorkingDirectory: util.StringOrNull(app.WorkingDirectory),
		LaunchParameters: util.StringOrNull(app.LaunchParameters),
		Platforms:        flattenPlatforms(ctx, app.Platforms, &resp.Diagnostics),
		InstanceFamilies: util.SetStringOrNull(ctx, app.InstanceFamilies, &resp.Diagnostics),
		AppBlockARN:      util.StringOrNull(app.AppBlockArn),
		Tags:             types.Map{},
		CreatedTime:      util.StringFromTime(app.CreatedTime),
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
