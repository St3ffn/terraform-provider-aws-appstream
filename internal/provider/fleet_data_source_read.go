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

func (ds *fleetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config fleetModel

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
			"Cannot read fleet because name must be set and known.",
		)
		return
	}

	name := config.Name.ValueString()

	out, err := ds.appstreamClient.DescribeFleets(ctx, &awsappstream.DescribeFleetsInput{
		Names: []string{name},
	})
	if err != nil {
		if isContextCanceled(ctx) {
			return
		}

		if isAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Fleet Not Found",
				fmt.Sprintf("No fleet named %q was found.", name),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream Fleet",
			fmt.Sprintf("Could not read fleet %q: %v", name, err),
		)
		return
	}

	if len(out.Fleets) == 0 {
		resp.Diagnostics.AddError(
			"AWS AppStream Fleet Not Found",
			fmt.Sprintf("No fleet named %q was found.", name),
		)
		return
	}

	fleet := out.Fleets[0]
	if fleet.Name == nil {
		resp.Diagnostics.AddError(
			"Unexpected AWS Response",
			fmt.Sprintf("Fleet %q was returned without required identifiers.", name),
		)
		return
	}

	state := &fleetModel{
		ID:                             types.StringValue(aws.ToString(fleet.Name)),
		Name:                           types.StringValue(aws.ToString(fleet.Name)),
		ImageName:                      stringOrNull(fleet.ImageName),
		ImageARN:                       stringOrNull(fleet.ImageArn),
		InstanceType:                   stringOrNull(fleet.InstanceType),
		FleetType:                      types.StringValue(string(fleet.FleetType)),
		ComputeCapacity:                flattenFleetComputeCapacity(ctx, fleet.ComputeCapacityStatus, &resp.Diagnostics),
		VPCConfig:                      flattenFleetVPCConfig(ctx, fleet.VpcConfig, &resp.Diagnostics),
		MaxUserDurationInSeconds:       int32OrNull(fleet.MaxUserDurationInSeconds),
		DisconnectTimeoutInSeconds:     int32OrNull(fleet.DisconnectTimeoutInSeconds),
		IdleDisconnectTimeoutInSeconds: int32OrNull(fleet.IdleDisconnectTimeoutInSeconds),
		Description:                    stringOrNull(fleet.Description),
		DisplayName:                    stringOrNull(fleet.DisplayName),
		EnableDefaultInternetAccess:    boolOrNull(fleet.EnableDefaultInternetAccess),
		DomainJoinInfo:                 flattenFleetDomainJoinInfo(ctx, fleet.DomainJoinInfo, &resp.Diagnostics),
		IAMRoleARN:                     stringOrNull(fleet.IamRoleArn),
		StreamView:                     types.StringValue(string(fleet.StreamView)),
		Platform:                       types.StringValue(string(fleet.Platform)),
		MaxConcurrentSessions:          int32OrNull(fleet.MaxConcurrentSessions),
		MaxSessionsPerInstance:         int32OrNull(fleet.MaxSessionsPerInstance),
		USBDeviceFilterStrings:         setStringOrNull(ctx, fleet.UsbDeviceFilterStrings, &resp.Diagnostics),
		SessionScriptS3Location:        flattenFleetSessionScriptS3Location(ctx, fleet.SessionScriptS3Location, &resp.Diagnostics),
		RootVolumeConfig:               flattenFleetRootVolumeConfig(ctx, fleet.RootVolumeConfig, &resp.Diagnostics),
		Tags:                           types.Map{},
		ARN:                            stringOrNull(fleet.Arn),
		CreatedTime:                    stringFromTime(fleet.CreatedTime),
		FleetErrors:                    flattenFleetErrors(ctx, fleet.FleetErrors, &resp.Diagnostics),
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
