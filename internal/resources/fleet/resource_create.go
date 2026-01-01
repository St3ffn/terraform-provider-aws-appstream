// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
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

	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create fleet because name must be known.",
		)
		return
	}

	name := plan.Name.ValueString()

	input := &awsappstream.CreateFleetInput{
		Name:         aws.String(name),
		InstanceType: aws.String(plan.InstanceType.ValueString()),
		FleetType:    awstypes.FleetType(plan.FleetType.ValueString()),
	}

	input.ImageName = util.StringPointerOrNil(plan.ImageName)
	input.ImageArn = util.StringPointerOrNil(plan.ImageARN)

	if !plan.ComputeCapacity.IsNull() && !plan.ComputeCapacity.IsUnknown() {
		input.ComputeCapacity = expandComputeCapacity(ctx, plan.ComputeCapacity, &resp.Diagnostics)
	}
	input.MaxConcurrentSessions = util.Int32PointerOrNil(plan.MaxConcurrentSessions)
	input.MaxSessionsPerInstance = util.Int32PointerOrNil(plan.MaxSessionsPerInstance)

	if !plan.VPCConfig.IsNull() && !plan.VPCConfig.IsUnknown() {
		input.VpcConfig = expandVPCConfig(ctx, plan.VPCConfig, &resp.Diagnostics)
	}

	input.MaxUserDurationInSeconds = util.Int32PointerOrNil(plan.MaxUserDurationInSeconds)
	input.DisconnectTimeoutInSeconds = util.Int32PointerOrNil(plan.DisconnectTimeoutInSeconds)
	input.IdleDisconnectTimeoutInSeconds = util.Int32PointerOrNil(plan.IdleDisconnectTimeoutInSeconds)

	input.Description = util.StringPointerOrNil(plan.Description)
	input.DisplayName = util.StringPointerOrNil(plan.DisplayName)
	input.EnableDefaultInternetAccess = util.BoolPointerOrNil(plan.EnableDefaultInternetAccess)
	input.IamRoleArn = util.StringPointerOrNil(plan.IAMRoleARN)

	if !plan.StreamView.IsNull() && !plan.StreamView.IsUnknown() {
		input.StreamView = awstypes.StreamView(plan.StreamView.ValueString())
	}

	if !plan.Platform.IsNull() && !plan.Platform.IsUnknown() {
		input.Platform = awstypes.PlatformType(plan.Platform.ValueString())
	}

	if !plan.DomainJoinInfo.IsNull() && !plan.DomainJoinInfo.IsUnknown() {
		input.DomainJoinInfo = expandDomainJoinInfo(ctx, plan.DomainJoinInfo, &resp.Diagnostics)
	}

	if !plan.SessionScriptS3Location.IsNull() && !plan.SessionScriptS3Location.IsUnknown() {
		input.SessionScriptS3Location = expandSessionScriptS3Location(
			ctx, plan.SessionScriptS3Location, &resp.Diagnostics,
		)
	}

	if !plan.RootVolumeConfig.IsNull() && !plan.RootVolumeConfig.IsUnknown() {
		input.RootVolumeConfig = expandRootVolumeConfig(
			ctx, plan.RootVolumeConfig, &resp.Diagnostics,
		)
	}

	input.UsbDeviceFilterStrings = util.ExpandStringSetOrNil(
		ctx, plan.USBDeviceFilterStrings, &resp.Diagnostics,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.appstreamClient.CreateFleet(ctx, input)
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsResourceAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Fleet Already Exists",
				fmt.Sprintf(
					"A fleet named %q already exists. To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> %q",
					name, name,
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Fleet",
			fmt.Sprintf("Could not create fleet %q: %v", name, err),
		)
		return
	}

	if out.Fleet != nil && out.Fleet.Arn != nil {
		_, tagDiags := r.tags.Apply(ctx, aws.ToString(out.Fleet.Arn), plan.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	newState, diags := r.readFleet(ctx, plan)
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
