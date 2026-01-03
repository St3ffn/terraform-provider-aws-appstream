// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	// Validate required attributes
	if plan.Name.IsNull() || plan.Name.IsUnknown() ||
		plan.IconS3Location.IsNull() || plan.IconS3Location.IsUnknown() ||
		plan.LaunchPath.IsNull() || plan.LaunchPath.IsUnknown() ||
		plan.Platforms.IsNull() || plan.Platforms.IsUnknown() ||
		plan.InstanceFamilies.IsNull() || plan.InstanceFamilies.IsUnknown() ||
		plan.AppBlockARN.IsNull() || plan.AppBlockARN.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create application because name, icon_s3_location, launch_path, platforms, instance_families, and app_block_arn must be known.",
		)
		return
	}

	name := plan.Name.ValueString()

	input := &awsappstream.CreateApplicationInput{
		Name:             aws.String(name),
		LaunchPath:       aws.String(plan.LaunchPath.ValueString()),
		AppBlockArn:      aws.String(plan.AppBlockARN.ValueString()),
		Platforms:        expandPlatforms(ctx, plan.Platforms, &resp.Diagnostics),
		InstanceFamilies: util.ExpandStringSetOrNil(ctx, plan.InstanceFamilies, &resp.Diagnostics),
	}

	input.DisplayName = util.StringPointerOrNil(plan.DisplayName)
	input.Description = util.StringPointerOrNil(plan.Description)
	input.WorkingDirectory = util.StringPointerOrNil(plan.WorkingDirectory)
	input.LaunchParameters = util.StringPointerOrNil(plan.LaunchParameters)

	input.IconS3Location = expandIconS3Location(
		ctx, plan.IconS3Location, &resp.Diagnostics,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	var out *awsappstream.CreateApplicationOutput
	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			var err error
			out, err = r.appstreamClient.CreateApplication(ctx, input)
			return err
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateApplication.html
		util.WithRetryOnFns(
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
			util.IsResourceNotFoundException,
		),
	)

	if err != nil {
		if util.IsResourceAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Application Already Exists",
				fmt.Sprintf(
					"An application named %q already exists. "+
						"To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> <application_arn>",
					name,
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Application",
			fmt.Sprintf("Could not create application %q: %v", name, err),
		)
		return
	}

	if out != nil && out.Application != nil && out.Application.Arn != nil {
		_, tagDiags := r.tags.Apply(ctx, aws.ToString(out.Application.Arn), plan.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	newState, diags := r.readApplication(ctx, aws.ToString(out.Application.Arn))
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
