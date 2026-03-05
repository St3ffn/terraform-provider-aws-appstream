// Copyright St3ffn 2025, 2026
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

// Update performs the full fleet reconciliation flow from state/plan diff.
// It builds UpdateFleet inputs (including AttributesToDelete), enforces stop/start
// preconditions based on update mode and update_behavior, retries transient API
// errors, applies tags, then converges desired_state and finally reads back state.
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

	diff := newResourceDiff(state, plan)
	name := plan.Name.ValueString()

	var input *awsappstream.UpdateFleetInput
	if diff.HasRemoteChanges() {
		input = &awsappstream.UpdateFleetInput{
			Name: aws.String(name),
		}

		var attrsToDelete []awstypes.FleetAttribute

		if diff.DisplayName.IsChanged() {
			util.OptionalStringUpdate(plan.DisplayName, state.DisplayName, func(v *string) {
				input.DisplayName = v
			})
		}

		if diff.Description.IsChanged() {
			util.OptionalStringUpdate(plan.Description, state.Description, func(v *string) {
				input.Description = v
			})
		}

		if diff.ImageName.IsChanged() {
			util.OptionalStringUpdate(plan.ImageName, state.ImageName, func(v *string) {
				input.ImageName = v
			})
		}

		if diff.ImageARN.IsChanged() {
			util.OptionalStringUpdate(plan.ImageARN, state.ImageARN, func(v *string) {
				input.ImageArn = v
			})
		}

		if diff.InstanceType.IsChanged() {
			util.OptionalStringUpdate(plan.InstanceType, state.InstanceType, func(v *string) {
				input.InstanceType = v
			})
		}

		if diff.IAMRoleARN.IsChanged() {
			if plan.IAMRoleARN.IsNull() {
				attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeIamRoleArn)
			} else {
				input.IamRoleArn = plan.IAMRoleARN.ValueStringPointer()
			}
		}

		if diff.StreamView.IsChanged() {
			// AWS does not support unsetting stream_view once configured.
			if !plan.StreamView.IsNull() {
				input.StreamView = awstypes.StreamView(plan.StreamView.ValueString())
			}
		}

		if diff.Platform.IsChanged() {
			// AWS does not support unsetting platform once configured.
			if !plan.Platform.IsNull() {
				input.Platform = awstypes.PlatformType(plan.Platform.ValueString())
			}
		}

		if diff.MaxUserDurationInSeconds.IsChanged() {
			// AWS does not support unsetting max_user_duration_in_seconds once configured.
			if !plan.MaxUserDurationInSeconds.IsNull() {
				input.MaxUserDurationInSeconds = plan.MaxUserDurationInSeconds.ValueInt32Pointer()
			}
		}

		if diff.DisconnectTimeoutInSeconds.IsChanged() {
			// AWS does not support unsetting disconnect_timeout_in_seconds once configured.
			if !plan.DisconnectTimeoutInSeconds.IsNull() {
				input.DisconnectTimeoutInSeconds = plan.DisconnectTimeoutInSeconds.ValueInt32Pointer()
			}
		}

		if diff.IdleDisconnectTimeoutInSeconds.IsChanged() {
			// AWS does not support unsetting idle_disconnect_timeout_in_seconds once configured.
			if !plan.IdleDisconnectTimeoutInSeconds.IsNull() {
				input.IdleDisconnectTimeoutInSeconds = plan.IdleDisconnectTimeoutInSeconds.ValueInt32Pointer()
			}
		}

		if diff.EnableDefaultInternetAccess.IsChanged() {
			// AWS does not support unsetting enable_default_internet_access once configured.
			if !plan.EnableDefaultInternetAccess.IsNull() {
				input.EnableDefaultInternetAccess = plan.EnableDefaultInternetAccess.ValueBoolPointer()
			}
		}

		if diff.DisableIMDSV1.IsChanged() {
			// AWS does not support unsetting disable_imdsv1 once configured.
			if !plan.DisableIMDSV1.IsNull() {
				input.DisableIMDSV1 = plan.DisableIMDSV1.ValueBoolPointer()
			}
		}

		if diff.ComputeCapacity.IsChanged() {
			// AWS does not support removing compute_capacity as a whole.
			if !plan.ComputeCapacity.IsNull() {
				input.ComputeCapacity = expandComputeCapacity(
					ctx,
					plan.ComputeCapacity,
					&resp.Diagnostics,
				)
			}
		}

		if diff.VPCConfig.IsChanged() {
			if plan.VPCConfig.IsNull() {
				attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeVpcConfiguration)
			} else {
				input.VpcConfig = expandVPCConfig(ctx, plan.VPCConfig, &resp.Diagnostics)
			}
		}

		if diff.DomainJoinInfo.IsChanged() {
			if plan.DomainJoinInfo.IsNull() {
				attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeDomainJoinInfo)
			} else {
				input.DomainJoinInfo = expandDomainJoinInfo(
					ctx,
					plan.DomainJoinInfo,
					&resp.Diagnostics,
				)
			}
		}

		if diff.USBDeviceFilterStrings.IsChanged() {
			if plan.USBDeviceFilterStrings.IsNull() {
				attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeUsbDeviceFilterStrings)
			} else {
				input.UsbDeviceFilterStrings = util.ExpandStringSetOrNil(
					ctx,
					plan.USBDeviceFilterStrings,
					&resp.Diagnostics,
				)
			}
		}

		if diff.SessionScriptS3Location.IsChanged() {
			if plan.SessionScriptS3Location.IsNull() {
				attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeSessionScriptS3Location)
			} else {
				input.SessionScriptS3Location = expandSessionScriptS3Location(
					ctx,
					plan.SessionScriptS3Location,
					&resp.Diagnostics,
				)
			}
		}

		if diff.MaxSessionsPerInstance.IsChanged() {
			if plan.MaxSessionsPerInstance.IsNull() {
				attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeMaxSessionsPerInstance)
			} else {
				input.MaxSessionsPerInstance = plan.MaxSessionsPerInstance.ValueInt32Pointer()
			}
		}

		if diff.RootVolumeConfig.IsChanged() {
			if plan.RootVolumeConfig.IsNull() {
				attrsToDelete = append(attrsToDelete, awstypes.FleetAttributeVolumeConfiguration)
			} else {
				input.RootVolumeConfig = expandRootVolumeConfig(
					ctx,
					plan.RootVolumeConfig,
					&resp.Diagnostics,
				)
			}
		}

		input.AttributesToDelete = attrsToDelete

		if resp.Diagnostics.HasError() {
			return
		}
	}

	restartAfter := false
	var err error
	mode := updateMode(plan, diff)
	updateBehavior := updateBehaviorFromPlan(plan.UpdateBehavior)
	desiredState := desiredStateFromPlan(plan.DesiredState)
	didRemoteUpdate := false
	arnForTags := ""
	switch {
	case !plan.ARN.IsNull() && !plan.ARN.IsUnknown():
		arnForTags = plan.ARN.ValueString()
	case !state.ARN.IsNull() && !state.ARN.IsUnknown():
		arnForTags = state.ARN.ValueString()
	}

	if diff.HasRemoteChanges() {
		switch mode {
		case fleetUpdateRequiresStop:
			switch updateBehavior {
			case updateBehaviorAutoStopStart:
				restartAfter, err = r.ensureStopped(ctx, name)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Updating AWS AppStream Fleet",
						fmt.Sprintf("Could not stop fleet %q for update: %v", name, err),
					)
					return
				}

			case updateBehaviorFailIfRunning:
				currentState, exists, err := r.currentFleetState(ctx, name)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Updating AWS AppStream Fleet",
						fmt.Sprintf("Could not inspect fleet %q state for update: %v", name, err),
					)
					return
				}

				if exists && currentState != awstypes.FleetStateStopped {
					resp.Diagnostics.AddError(
						"Error Updating AWS AppStream Fleet",
						fmt.Sprintf(
							"Could not update fleet %q: update requires the fleet to be stopped, but current state is %q and update_behavior is %q.",
							name,
							string(currentState),
							updateBehaviorFailIfRunning,
						),
					)
					return
				}

			default:
				resp.Diagnostics.AddError(
					"Error Updating AWS AppStream Fleet",
					fmt.Sprintf("Could not update fleet %q: unsupported update_behavior %q", name, updateBehavior),
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

		var out *awsappstream.UpdateFleetOutput
		out, err = util.RetryOnValue(
			ctx,
			func(ctx context.Context) (*awsappstream.UpdateFleetOutput, error) {
				return r.appstreamClient.UpdateFleet(ctx, input)
			},
			util.WithTimeout(updateRetryTimeout),
			util.WithInitBackoff(updateRetryInitBackoff),
			util.WithMaxBackoff(updateRetryMaxBackoff),
			// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_UpdateFleet.html
			util.WithRetryOnFns(
				util.IsConcurrentModificationException,
				util.IsOperationNotPermittedException,
				util.IsResourceNotAvailableException,
				util.IsResourceInUseException,
				util.IsInvalidRoleException,
			),
		)

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

		didRemoteUpdate = true
		if out.Fleet != nil && out.Fleet.Arn != nil {
			arnForTags = aws.ToString(out.Fleet.Arn)
		}
	}

	if diff.RequiresTagApply() {
		if arnForTags == "" {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Fleet",
				fmt.Sprintf("Could not apply tags for fleet %q: missing ARN in state and plan", name),
			)
			return
		}
		_, tagDiags := r.tags.Apply(ctx, arnForTags, plan.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	switch {
	case didRemoteUpdate && mode == fleetUpdateRequiresStop && desiredState == desiredStateRunning:
		err = r.ensureRunning(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Fleet",
				fmt.Sprintf("Could not start fleet %q after successful update: %v", name, err),
			)
			return
		}
	case didRemoteUpdate && mode == fleetUpdateRequiresStop && desiredState == desiredStateInherit && restartAfter:
		err = r.ensureRunning(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Fleet",
				fmt.Sprintf("Could not start fleet %q after successful update: %v", name, err),
			)
			return
		}
	case didRemoteUpdate && mode == fleetUpdateAllowedRunning && desiredState == desiredStateRunning:
		err = r.ensureRunning(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Fleet",
				fmt.Sprintf("Could not start fleet %q after successful update: %v", name, err),
			)
			return
		}
	case didRemoteUpdate && mode == fleetUpdateAllowedRunning && desiredState == desiredStateStopped:
		_, err = r.ensureStopped(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Fleet",
				fmt.Sprintf("Could not stop fleet %q after successful update: %v", name, err),
			)
			return
		}
	case !didRemoteUpdate && desiredState == desiredStateRunning:
		err = r.ensureRunning(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Fleet",
				fmt.Sprintf("Could not start fleet %q after successful update: %v", name, err),
			)
			return
		}
	case !didRemoteUpdate && desiredState == desiredStateStopped:
		_, err = r.ensureStopped(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Fleet",
				fmt.Sprintf("Could not stop fleet %q after successful update: %v", name, err),
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

func (r *resource) currentFleetState(ctx context.Context, name string) (state awstypes.FleetState, exists bool, err error) {
	out, err := r.appstreamClient.DescribeFleets(ctx, &awsappstream.DescribeFleetsInput{
		Names: []string{name},
	})
	if err != nil {
		if util.IsAppStreamNotFound(err) {
			return "", false, nil
		}
		return "", false, err
	}

	if len(out.Fleets) == 0 {
		return "", false, nil
	}

	return out.Fleets[0].State, true, nil
}

func (r *resource) ensureStopped(ctx context.Context, name string) (restartAfter bool, err error) {
	stateCaptured := false

	err = util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			state, exists, err := r.currentFleetState(ctx, name)
			if err != nil {
				return err
			}

			if !exists {
				return nil
			}

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
		util.WithTimeout(stopRetryTimeout),
		util.WithInitBackoff(stopRetryInitBackoff),
		util.WithMaxBackoff(stopRetryMaxBackoff),
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
			state, exists, err := r.currentFleetState(ctx, name)
			if err != nil {
				return err
			}

			if !exists {
				return nil
			}

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
		util.WithTimeout(startRetryTimeout),
		util.WithInitBackoff(startRetryInitBackoff),
		util.WithMaxBackoff(startRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeFleets.html
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_StartFleet.html
		util.WithRetryOnFns(
			isUnexpectedFleetStateError,
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
		),
	)
}
