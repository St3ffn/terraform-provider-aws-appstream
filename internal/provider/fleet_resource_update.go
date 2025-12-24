// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *fleetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan fleetModel
	var state fleetModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update fleet because name must be known.",
		)
		return
	}

	name := plan.Name.ValueString()

	// guard against unexpected identity drift
	if !state.Name.IsNull() && !state.Name.IsUnknown() {
		if state.Name.ValueString() != name {
			resp.Diagnostics.AddError(
				"Unexpected Update Request",
				"Fleet identity (name) changed during update. This should trigger replacement. Please report this issue.",
			)
			return
		}
	}

	input := &awsappstream.UpdateFleetInput{
		Name: aws.String(name),
	}

	var attrsToDelete []awstypes.FleetAttribute

	optionalStringUpdate(plan.DisplayName, state.DisplayName, func(v *string) {
		input.DisplayName = v
	})

	optionalStringUpdate(plan.Description, state.Description, func(v *string) {
		input.Description = v
	})

	optionalStringUpdate(plan.ImageName, state.ImageName, func(v *string) {
		input.ImageName = v
	})

	optionalStringUpdate(plan.ImageARN, state.ImageARN, func(v *string) {
		input.ImageArn = v
	})

	optionalStringUpdate(plan.InstanceType, state.InstanceType, func(v *string) {
		input.InstanceType = v
	})

	optionalStringUpdate(plan.IAMRoleARN, state.IAMRoleARN, func(v *string) {
		input.IamRoleArn = v
	})

	if plan.StreamView.IsNull() {
		// no delete support
	} else if !plan.StreamView.IsUnknown() {
		input.StreamView = awstypes.StreamView(plan.StreamView.ValueString())
	}

	if plan.Platform.IsNull() {
		// no delete support
	} else if !plan.Platform.IsUnknown() {
		input.Platform = awstypes.PlatformType(plan.Platform.ValueString())
	}

	if plan.MaxUserDurationInSeconds.IsNull() {
		// no delete support
	} else if !plan.MaxUserDurationInSeconds.IsUnknown() {
		input.MaxUserDurationInSeconds = plan.MaxUserDurationInSeconds.ValueInt32Pointer()
	}

	if plan.DisconnectTimeoutInSeconds.IsNull() {
		// no delete support
	} else if !plan.DisconnectTimeoutInSeconds.IsUnknown() {
		input.DisconnectTimeoutInSeconds = plan.DisconnectTimeoutInSeconds.ValueInt32Pointer()
	}

	if plan.IdleDisconnectTimeoutInSeconds.IsNull() {
		// no delete support
	} else if !plan.IdleDisconnectTimeoutInSeconds.IsUnknown() {
		input.IdleDisconnectTimeoutInSeconds = plan.IdleDisconnectTimeoutInSeconds.ValueInt32Pointer()
	}

	if plan.EnableDefaultInternetAccess.IsNull() {
		// no delete support
	} else if !plan.EnableDefaultInternetAccess.IsUnknown() {
		input.EnableDefaultInternetAccess = plan.EnableDefaultInternetAccess.ValueBoolPointer()
	}

	if plan.ComputeCapacity.IsNull() {
		// no delete support
	} else if !plan.ComputeCapacity.IsUnknown() {
		input.ComputeCapacity = expandFleetComputeCapacity(
			ctx,
			plan.ComputeCapacity,
			&resp.Diagnostics,
		)
	}

	if plan.VPCConfig.IsNull() {
		if !plan.VPCConfig.IsUnknown() {
			attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeVpcConfiguration)
		}
	} else if !plan.VPCConfig.IsUnknown() {
		input.VpcConfig = expandFleetVPCConfig(ctx, plan.VPCConfig, &resp.Diagnostics)
	}

	if plan.DomainJoinInfo.IsNull() {
		if !plan.DomainJoinInfo.IsUnknown() {
			attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeDomainJoinInfo)
		}
	} else if !plan.DomainJoinInfo.IsUnknown() {
		input.DomainJoinInfo = expandFleetDomainJoinInfo(
			ctx,
			plan.DomainJoinInfo,
			&resp.Diagnostics,
		)
	}

	if plan.USBDeviceFilterStrings.IsNull() {
		if !plan.USBDeviceFilterStrings.IsUnknown() {
			attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeUsbDeviceFilterStrings)
		}
	} else if !plan.USBDeviceFilterStrings.IsUnknown() {
		input.UsbDeviceFilterStrings = expandStringSetOrNil(
			ctx,
			plan.USBDeviceFilterStrings,
			&resp.Diagnostics,
		)
	}

	if plan.SessionScriptS3Location.IsNull() {
		if !plan.SessionScriptS3Location.IsUnknown() {
			attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeSessionScriptS3Location)
		}
	} else if !plan.SessionScriptS3Location.IsUnknown() {
		input.SessionScriptS3Location = expandFleetSessionScriptS3Location(
			ctx,
			plan.SessionScriptS3Location,
			&resp.Diagnostics,
		)
	}

	if plan.MaxSessionsPerInstance.IsNull() {
		if !plan.MaxSessionsPerInstance.IsUnknown() {
			attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeMaxSessionsPerInstance)
		}
	} else if !plan.MaxSessionsPerInstance.IsUnknown() {
		input.MaxSessionsPerInstance = plan.MaxSessionsPerInstance.ValueInt32Pointer()
	}

	if plan.RootVolumeConfig.IsNull() {
		if !plan.RootVolumeConfig.IsUnknown() {
			attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeVolumeConfiguration)
		}
	} else if !plan.RootVolumeConfig.IsUnknown() {
		input.RootVolumeConfig = expandFleetRootVolumeConfig(
			ctx,
			plan.RootVolumeConfig,
			&resp.Diagnostics,
		)
	}

	input.AttributesToDelete = attrsToDelete

	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.appstreamClient.UpdateFleet(ctx, input)
	if err != nil {
		if isContextCanceled(err) {
			return
		}

		if isAppStreamNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Updating AWS AppStream Fleet",
			fmt.Sprintf("Could not update fleet %q: %v", name, err),
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

	newState, diags := r.readFleet(ctx, name)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
