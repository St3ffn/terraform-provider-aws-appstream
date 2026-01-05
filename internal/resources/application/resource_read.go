// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

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

	arn := state.ID.ValueString()

	newState, diags := r.readApplication(ctx, arn)
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

func (r *resource) readApplication(ctx context.Context, arn string) (*model, diag.Diagnostics) {
	var diags diag.Diagnostics

	out, err := r.appstreamClient.DescribeApplications(ctx, &awsappstream.DescribeApplicationsInput{
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

	state := &model{
		ID:               types.StringValue(aws.ToString(app.Arn)),
		Name:             types.StringValue(aws.ToString(app.Name)),
		DisplayName:      util.StringOrNull(app.DisplayName),
		Description:      util.StringOrNull(app.Description),
		IconS3Location:   flattenIconS3Location(ctx, app.IconS3Location, &diags),
		LaunchPath:       util.StringOrNull(app.LaunchPath),
		WorkingDirectory: util.StringOrNull(app.WorkingDirectory),
		LaunchParameters: util.StringOrNull(app.LaunchParameters),
		Platforms:        flattenPlatforms(ctx, app.Platforms, &diags),
		InstanceFamilies: util.SetStringOrNull(ctx, app.InstanceFamilies, &diags),
		AppBlockARN:      util.StringOrNull(app.AppBlockArn),
		Tags:             types.MapNull(types.StringType),
		ARN:              util.StringOrNull(app.Arn),
		CreatedTime:      util.StringFromTime(app.CreatedTime),
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
