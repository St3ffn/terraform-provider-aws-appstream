// Copyright St3ffn 2025, 2026
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
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/tags"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Read(ctx context.Context, req tfresource.ReadRequest, resp *tfresource.ReadResponse) {
	var state resourceModel

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

	newState, diags := r.readFleet(ctx, state)
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

func (r *resource) readFleet(ctx context.Context, prior resourceModel) (*resourceModel, diag.Diagnostics) {
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

	state := &resourceModel{
		ID:                             types.StringValue(aws.ToString(fleet.Name)),
		Name:                           types.StringValue(aws.ToString(fleet.Name)),
		ImageName:                      util.StringOrNull(fleet.ImageName),
		ImageARN:                       util.StringOrNull(fleet.ImageArn),
		InstanceType:                   util.StringOrNull(fleet.InstanceType),
		FleetType:                      types.StringValue(string(fleet.FleetType)),
		ComputeCapacity:                flattenComputeCapacity(ctx, fleet.ComputeCapacityStatus, &diags),
		VPCConfig:                      flattenVPCConfig(ctx, fleet.VpcConfig, &diags),
		MaxUserDurationInSeconds:       util.Int32OrNull(fleet.MaxUserDurationInSeconds),
		DisconnectTimeoutInSeconds:     util.Int32OrNull(fleet.DisconnectTimeoutInSeconds),
		IdleDisconnectTimeoutInSeconds: util.Int32OrNull(fleet.IdleDisconnectTimeoutInSeconds),
		Description:                    util.StringOrNull(fleet.Description),
		DisplayName:                    util.StringOrNull(fleet.DisplayName),
		EnableDefaultInternetAccess:    util.BoolOrNull(fleet.EnableDefaultInternetAccess),
		DomainJoinInfo:                 flattenDomainJoinInfo(ctx, fleet.DomainJoinInfo, &diags),
		IAMRoleARN:                     util.StringOrNull(fleet.IamRoleArn),
		StreamView:                     util.FlattenStateOwnedString(prior.StreamView, aws.String(string(fleet.StreamView))),
		Platform:                       util.FlattenStateOwnedString(prior.Platform, aws.String(string(fleet.Platform))),
		MaxConcurrentSessions:          util.Int32OrNull(fleet.MaxConcurrentSessions),
		MaxSessionsPerInstance:         util.Int32OrNull(fleet.MaxSessionsPerInstance),
		USBDeviceFilterStrings:         util.SetStringOrNull(ctx, fleet.UsbDeviceFilterStrings, &diags),
		SessionScriptS3Location:        flattenSessionScriptS3Location(ctx, fleet.SessionScriptS3Location, &diags),
		RootVolumeConfig:               flattenRootVolumeConfig(ctx, fleet.RootVolumeConfig, &diags),
		Tags:                           types.MapNull(types.StringType),
		ARN:                            util.StringOrNull(fleet.Arn),
		CreatedTime:                    util.StringFromTime(fleet.CreatedTime),
		State:                          types.StringValue(string(fleet.State)),
		FleetErrors:                    flattenFleetErrors(ctx, fleet.FleetErrors, &diags),
	}

	if !state.ARN.IsNull() {
		allTags, allTagDiags := r.tags.ReadAll(ctx, state.ARN.ValueString())
		diags.Append(allTagDiags...)
		state.TagsAll = allTags

		resourceTags, resourceTagDiags := tags.ResourceTags(ctx, prior.Tags, allTags)
		diags.Append(resourceTagDiags...)
		state.Tags = resourceTags
	}

	if diags.HasError() {
		return nil, diags
	}

	return state, diags
}
