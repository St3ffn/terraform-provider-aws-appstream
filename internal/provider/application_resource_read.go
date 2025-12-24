// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state applicationModel

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

	arn := state.ID.ValueString()

	newState, diags := r.readApplication(ctx, arn)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if isContextCanceled(ctx.Err()) {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *applicationResource) readApplication(ctx context.Context, arn string) (*applicationModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	out, err := r.appstreamClient.DescribeApplications(ctx, &awsappstream.DescribeApplicationsInput{
		Arns: []string{arn},
	})
	if err != nil {
		if isContextCanceled(err) {
			return nil, diags
		}

		if isAppStreamNotFound(err) {
			return nil, diags
		}

		diags.AddError(
			"Error Reading AWS AppStream Application",
			fmt.Sprintf("Could not read application %q: %v", arn, err),
		)
		return nil, diags
	}

	if len(out.Applications) == 0 {
		return nil, diags
	}

	app := out.Applications[0]
	if app.Arn == nil || app.Name == nil {
		return nil, diags
	}

	state := &applicationModel{
		ID:               types.StringValue(aws.ToString(app.Arn)),
		Name:             types.StringValue(aws.ToString(app.Name)),
		DisplayName:      stringOrNull(app.DisplayName),
		Description:      stringOrNull(app.Description),
		IconS3Location:   flattenApplicationIconS3Location(ctx, app.IconS3Location, &diags),
		LaunchPath:       stringOrNull(app.LaunchPath),
		WorkingDirectory: stringOrNull(app.WorkingDirectory),
		LaunchParameters: stringOrNull(app.LaunchParameters),
		Platforms:        flattenApplicationPlatforms(ctx, app.Platforms, &diags),
		InstanceFamilies: setStringOrNull(ctx, app.InstanceFamilies, &diags),
		AppBlockARN:      stringOrNull(app.AppBlockArn),
		Tags:             types.Map{},
		ARN:              stringOrNull(app.Arn),
		CreatedTime:      stringFromTime(app.CreatedTime),
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
