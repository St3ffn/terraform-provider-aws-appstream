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

func (r *resource) Update(ctx context.Context, req tfresource.UpdateRequest, resp *tfresource.UpdateResponse) {
	var plan resourceModel
	var state resourceModel

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

	input := &awsappstream.UpdateFleetInput{
		Name: aws.String(name),
	}

	var attrsToDelete []awstypes.FleetAttribute

	util.OptionalStringUpdate(plan.DisplayName, state.DisplayName, func(v *string) {
		input.DisplayName = v
	})

	util.OptionalStringUpdate(plan.Description, state.Description, func(v *string) {
		input.Description = v
	})

	util.OptionalStringUpdate(plan.ImageName, state.ImageName, func(v *string) {
		input.ImageName = v
	})

	util.OptionalStringUpdate(plan.ImageARN, state.ImageARN, func(v *string) {
		input.ImageArn = v
	})

	util.OptionalStringUpdate(plan.InstanceType, state.InstanceType, func(v *string) {
		input.InstanceType = v
	})

	util.OptionalStringUpdate(plan.IAMRoleARN, state.IAMRoleARN, func(v *string) {
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
		input.ComputeCapacity = expandComputeCapacity(
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
		input.VpcConfig = expandVPCConfig(ctx, plan.VPCConfig, &resp.Diagnostics)
	}

	if plan.DomainJoinInfo.IsNull() {
		if !plan.DomainJoinInfo.IsUnknown() {
			attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeDomainJoinInfo)
		}
	} else if !plan.DomainJoinInfo.IsUnknown() {
		input.DomainJoinInfo = expandDomainJoinInfo(
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
		input.UsbDeviceFilterStrings = util.ExpandStringSetOrNil(
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
		input.SessionScriptS3Location = expandSessionScriptS3Location(
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
		input.RootVolumeConfig = expandRootVolumeConfig(
			ctx,
			plan.RootVolumeConfig,
			&resp.Diagnostics,
		)
	}

	input.AttributesToDelete = attrsToDelete

	if resp.Diagnostics.HasError() {
		return
	}

	restartAfter := false
	var err error
	mode := updateMode(state, plan)

	switch mode {
	case fleetUpdateRequiresStop:
		restartAfter, err = r.ensureStopped(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Fleet",
				fmt.Sprintf("Could not stop fleet %q for update: %v", name, err),
			)
			return
		}
	case fleetUpdateAllowedRunning:
		// proceed normally
	case fleetUpdateForbidden:
		resp.Diagnostics.AddError(
			"Error Updating AWS AppStream Fleet",
			fmt.Sprintf("Could not update fleet %q: fleet cannot be updated in its current state", name),
		)
		return
	}

	out, err := r.appstreamClient.UpdateFleet(ctx, input)
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
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

	if mode == fleetUpdateRequiresStop && restartAfter {
		err = r.ensureRunning(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Fleet",
				fmt.Sprintf("Could not start fleet %q after successful update: %v", name, err),
			)
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

func wasRunning(state awstypes.FleetState) bool {
	return state == awstypes.FleetStateRunning || state == awstypes.FleetStateStarting
}

func (r *resource) ensureStopped(ctx context.Context, name string) (restartAfter bool, err error) {
	stateCaptured := false

	err = util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			out, err := r.appstreamClient.DescribeFleets(ctx, &awsappstream.DescribeFleetsInput{
				Names: []string{name},
			})
			if err != nil {
				if util.IsAppStreamNotFound(err) {
					return nil
				}
				return err
			}

			if len(out.Fleets) == 0 {
				return nil
			}

			state := out.Fleets[0].State

			if !stateCaptured {
				restartAfter = wasRunning(state)
				stateCaptured = true
			}

			switch state {
			case awstypes.FleetStateRunning:
				// stoppable state
				_, err = r.appstreamClient.StopFleet(ctx, &awsappstream.StopFleetInput{
					Name: aws.String(name),
				})
				if err != nil {
					if util.IsAppStreamNotFound(err) {
						return nil
					}
					return err
				}
				// retry as we just stopped the fleet
				return fmt.Errorf("%w: current=%s", errUnexpectedFleetState, state)

			case awstypes.FleetStateStopped:
				// updatable state
				return nil

			default:
				return fmt.Errorf("%w: current=%s", errUnexpectedFleetState, state)
			}
		},
		util.WithTimeout(updateRetryTimeout),
		util.WithInitBackoff(updateRetryInitBackoff),
		util.WithMaxBackoff(updateRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeFleets.html
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_StopFleet.html
		util.WithRetryOnFns(
			isUnexpectedFleetStateError,
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
		),
	)
	return
}

func (r *resource) ensureRunning(ctx context.Context, name string) error {
	return util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			out, err := r.appstreamClient.DescribeFleets(ctx, &awsappstream.DescribeFleetsInput{
				Names: []string{name},
			})
			if err != nil {
				if util.IsAppStreamNotFound(err) {
					return nil
				}
				return err
			}

			if len(out.Fleets) == 0 {
				return nil
			}

			state := out.Fleets[0].State

			switch state {
			case awstypes.FleetStateStopped:
				// startable state
				_, err = r.appstreamClient.StartFleet(ctx, &awsappstream.StartFleetInput{
					Name: aws.String(name),
				})
				if err != nil {
					if util.IsAppStreamNotFound(err) {
						return nil
					}
					return err
				}
				// retry as we just started the fleet
				return fmt.Errorf("%w: current=%s", errUnexpectedFleetState, state)

			case awstypes.FleetStateRunning:
				// desired state
				return nil

			default:
				return fmt.Errorf("%w: current=%s", errUnexpectedFleetState, state)
			}
		},
		util.WithTimeout(updateRetryTimeout),
		util.WithInitBackoff(updateRetryInitBackoff),
		util.WithMaxBackoff(updateRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeFleets.html
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_StartFleet.html
		util.WithRetryOnFns(
			isUnexpectedFleetStateError,
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
		),
	)
}
