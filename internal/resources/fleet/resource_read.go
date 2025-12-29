// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet

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

	if state.Name.IsNull() || state.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Required attribute name is missing from state. "+
				"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the tfresource.",
		)
		return
	}

	newState, diags := r.readFleet(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if util.IsContextCanceled(ctx.Err()) {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *resource) readFleet(ctx context.Context, prior model) (*model, diag.Diagnostics) {
	var diags diag.Diagnostics

	name := prior.Name.ValueString()

	out, err := r.appstreamClient.DescribeFleets(ctx, &awsappstream.DescribeFleetsInput{
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
			"Error Reading AWS AppStream Fleet",
			fmt.Sprintf("Could not read fleet %q: %v", name, err),
		)
		return nil, diags
	}

	if len(out.Fleets) == 0 {
		return nil, diags
	}

	fleet := out.Fleets[0]
	if fleet.Name == nil {
		return nil, diags
	}

	state := &model{
		ID:                             types.StringValue(aws.ToString(fleet.Name)),
		Name:                           types.StringValue(aws.ToString(fleet.Name)),
		ImageName:                      util.FlattenOwnedString(prior.ImageName, fleet.ImageName),
		ImageARN:                       util.FlattenOwnedString(prior.ImageARN, fleet.ImageArn),
		InstanceType:                   util.StringOrNull(fleet.InstanceType),
		FleetType:                      types.StringValue(string(fleet.FleetType)),
		ComputeCapacity:                flattenComputeCapacity(ctx, fleet.ComputeCapacityStatus, &diags),
		VPCConfig:                      flattenVPCConfig(ctx, fleet.VpcConfig, &diags),
		MaxUserDurationInSeconds:       util.FlattenOwnedInt32(prior.MaxUserDurationInSeconds, fleet.MaxUserDurationInSeconds),
		DisconnectTimeoutInSeconds:     util.FlattenOwnedInt32(prior.DisconnectTimeoutInSeconds, fleet.DisconnectTimeoutInSeconds),
		IdleDisconnectTimeoutInSeconds: util.FlattenOwnedInt32(prior.IdleDisconnectTimeoutInSeconds, fleet.IdleDisconnectTimeoutInSeconds),
		Description:                    util.StringOrNull(fleet.Description),
		DisplayName:                    util.StringOrNull(fleet.DisplayName),
		EnableDefaultInternetAccess:    util.FlattenOwnedBool(prior.EnableDefaultInternetAccess, fleet.EnableDefaultInternetAccess),
		DomainJoinInfo:                 flattenDomainJoinInfo(ctx, fleet.DomainJoinInfo, &diags),
		IAMRoleARN:                     util.StringOrNull(fleet.IamRoleArn),
		StreamView:                     util.FlattenOwnedString(prior.StreamView, aws.String(string(fleet.StreamView))),
		Platform:                       util.FlattenOwnedString(prior.Platform, aws.String(string(fleet.Platform))),
		MaxConcurrentSessions:          util.Int32OrNull(fleet.MaxConcurrentSessions),
		MaxSessionsPerInstance:         util.Int32OrNull(fleet.MaxSessionsPerInstance),
		USBDeviceFilterStrings:         util.SetStringOrNull(ctx, fleet.UsbDeviceFilterStrings, &diags),
		SessionScriptS3Location:        flattenSessionScriptS3Location(ctx, fleet.SessionScriptS3Location, &diags),
		RootVolumeConfig:               flattenRootVolumeConfig(ctx, fleet.RootVolumeConfig, &diags),
		Tags:                           types.Map{},
		ARN:                            util.StringOrNull(fleet.Arn),
		CreatedTime:                    util.StringFromTime(fleet.CreatedTime),
		FleetErrors:                    flattenFleetErrors(ctx, fleet.FleetErrors, &diags),
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
