// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (ds *applicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config applicationModel

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
		if isContextCanceled(err) {
			return
		}

		if isAppStreamNotFound(err) {
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

	state := &applicationModel{
		ID:               types.StringValue(aws.ToString(app.Arn)),
		ARN:              types.StringValue(aws.ToString(app.Arn)),
		Name:             types.StringValue(aws.ToString(app.Name)),
		DisplayName:      stringOrNull(app.DisplayName),
		Description:      stringOrNull(app.Description),
		IconS3Location:   flattenApplicationIconS3Location(ctx, app.IconS3Location, &resp.Diagnostics),
		LaunchPath:       stringOrNull(app.LaunchPath),
		WorkingDirectory: stringOrNull(app.WorkingDirectory),
		LaunchParameters: stringOrNull(app.LaunchParameters),
		Platforms:        flattenApplicationPlatforms(ctx, app.Platforms, &resp.Diagnostics),
		InstanceFamilies: setStringOrNull(ctx, app.InstanceFamilies, &resp.Diagnostics),
		AppBlockARN:      stringOrNull(app.AppBlockArn),
		Tags:             types.Map{},
		CreatedTime:      stringFromTime(app.CreatedTime),
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
