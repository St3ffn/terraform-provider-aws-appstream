// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet

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
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
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

	state := &model{
		ID:                             types.StringValue(aws.ToString(fleet.Name)),
		Name:                           types.StringValue(aws.ToString(fleet.Name)),
		ImageName:                      util.StringOrNull(fleet.ImageName),
		ImageARN:                       util.StringOrNull(fleet.ImageArn),
		InstanceType:                   util.StringOrNull(fleet.InstanceType),
		FleetType:                      types.StringValue(string(fleet.FleetType)),
		ComputeCapacity:                flattenComputeCapacity(ctx, fleet.ComputeCapacityStatus, &resp.Diagnostics),
		VPCConfig:                      flattenVPCConfig(ctx, fleet.VpcConfig, &resp.Diagnostics),
		MaxUserDurationInSeconds:       util.Int32OrNull(fleet.MaxUserDurationInSeconds),
		DisconnectTimeoutInSeconds:     util.Int32OrNull(fleet.DisconnectTimeoutInSeconds),
		IdleDisconnectTimeoutInSeconds: util.Int32OrNull(fleet.IdleDisconnectTimeoutInSeconds),
		Description:                    util.StringOrNull(fleet.Description),
		DisplayName:                    util.StringOrNull(fleet.DisplayName),
		EnableDefaultInternetAccess:    util.BoolOrNull(fleet.EnableDefaultInternetAccess),
		DomainJoinInfo:                 flattenDomainJoinInfo(ctx, fleet.DomainJoinInfo, &resp.Diagnostics),
		IAMRoleARN:                     util.StringOrNull(fleet.IamRoleArn),
		StreamView:                     types.StringValue(string(fleet.StreamView)),
		Platform:                       types.StringValue(string(fleet.Platform)),
		MaxConcurrentSessions:          util.Int32OrNull(fleet.MaxConcurrentSessions),
		MaxSessionsPerInstance:         util.Int32OrNull(fleet.MaxSessionsPerInstance),
		USBDeviceFilterStrings:         util.SetStringOrNull(ctx, fleet.UsbDeviceFilterStrings, &resp.Diagnostics),
		SessionScriptS3Location:        flattenSessionScriptS3Location(ctx, fleet.SessionScriptS3Location, &resp.Diagnostics),
		RootVolumeConfig:               flattenRootVolumeConfig(ctx, fleet.RootVolumeConfig, &resp.Diagnostics),
		Tags:                           types.Map{},
		ARN:                            util.StringOrNull(fleet.Arn),
		CreatedTime:                    util.StringFromTime(fleet.CreatedTime),
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
