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

func (r *fleetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state fleetModel

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
				"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
		)
		return
	}

	name := state.Name.ValueString()

	newState, diags := r.readFleet(ctx, name)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if isContextCanceled(ctx) {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *fleetResource) readFleet(ctx context.Context, name string) (*fleetModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	out, err := r.appstreamClient.DescribeFleets(ctx, &awsappstream.DescribeFleetsInput{
		Names: []string{name},
	})
	if err != nil {
		if isContextCanceled(ctx) {
			return nil, diags
		}

		if isAppStreamNotFound(err) {
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

	state := &fleetModel{
		ID:                             types.StringValue(aws.ToString(fleet.Name)),
		Name:                           types.StringValue(aws.ToString(fleet.Name)),
		ImageName:                      stringOrNull(fleet.ImageName),
		ImageARN:                       stringOrNull(fleet.ImageArn),
		InstanceType:                   stringOrNull(fleet.InstanceType),
		FleetType:                      types.StringValue(string(fleet.FleetType)),
		ComputeCapacity:                flattenFleetComputeCapacity(ctx, fleet.ComputeCapacityStatus, &diags),
		VPCConfig:                      flattenFleetVPCConfig(ctx, fleet.VpcConfig, &diags),
		MaxUserDurationInSeconds:       int32OrNull(fleet.MaxUserDurationInSeconds),
		DisconnectTimeoutInSeconds:     int32OrNull(fleet.DisconnectTimeoutInSeconds),
		IdleDisconnectTimeoutInSeconds: int32OrNull(fleet.IdleDisconnectTimeoutInSeconds),
		Description:                    stringOrNull(fleet.Description),
		DisplayName:                    stringOrNull(fleet.DisplayName),
		EnableDefaultInternetAccess:    boolOrNull(fleet.EnableDefaultInternetAccess),
		DomainJoinInfo:                 flattenFleetDomainJoinInfo(ctx, fleet.DomainJoinInfo, &diags),
		IAMRoleARN:                     stringOrNull(fleet.IamRoleArn),
		StreamView:                     types.StringValue(string(fleet.StreamView)),
		Platform:                       types.StringValue(string(fleet.Platform)),
		MaxConcurrentSessions:          int32OrNull(fleet.MaxConcurrentSessions),
		MaxSessionsPerInstance:         int32OrNull(fleet.MaxSessionsPerInstance),
		USBDeviceFilterStrings:         setStringOrNull(ctx, fleet.UsbDeviceFilterStrings, &diags),
		SessionScriptS3Location:        flattenFleetSessionScriptS3Location(ctx, fleet.SessionScriptS3Location, &diags),
		RootVolumeConfig:               flattenFleetRootVolumeConfig(ctx, fleet.RootVolumeConfig, &diags),
		Tags:                           types.Map{},
		ARN:                            stringOrNull(fleet.Arn),
		CreatedTime:                    stringFromTime(fleet.CreatedTime),
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
