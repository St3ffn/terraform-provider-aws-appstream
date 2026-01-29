// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import (
	"context"
	"fmt"
	"strings"

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
			"Cannot update app block builder because name must be known.",
		)
		return
	}

	name := plan.Name.ValueString()

	input := &awsappstream.UpdateAppBlockBuilderInput{
		Name: aws.String(name),
	}

	var attrsToDelete []awstypes.AppBlockBuilderAttribute

	util.OptionalStringUpdate(plan.Description, state.Description, func(v *string) {
		input.Description = v
	})

	util.OptionalStringUpdate(plan.DisplayName, state.DisplayName, func(v *string) {
		input.DisplayName = v
	})

	input.EnableDefaultInternetAccess = util.BoolPointerOrNil(plan.EnableDefaultInternetAccess)

	if plan.InstanceType.IsNull() {
		// no delete support
	} else if !plan.InstanceType.IsUnknown() {
		input.InstanceType = plan.InstanceType.ValueStringPointer()
	}

	if plan.Platform.IsNull() {
		// no delete support
	} else if !plan.Platform.IsUnknown() {
		input.Platform = awstypes.PlatformType(plan.Platform.ValueString())
	}

	if plan.IAMRoleARN.IsNull() {
		if !plan.IAMRoleARN.IsUnknown() {
			attrsToDelete = append(attrsToDelete, awstypes.AppBlockBuilderAttributeIamRoleArn)
		}
	} else if !plan.IAMRoleARN.IsUnknown() {
		input.IamRoleArn = plan.IAMRoleARN.ValueStringPointer()
	}

	if plan.AccessEndpoints.IsNull() {
		if !plan.AccessEndpoints.IsUnknown() {
			attrsToDelete = append(attrsToDelete, awstypes.AppBlockBuilderAttributeAccessEndpoints)
		}
	} else if !plan.AccessEndpoints.IsUnknown() {
		input.AccessEndpoints = expandAccessEndpoints(ctx, plan.AccessEndpoints, &resp.Diagnostics)
	}

	if plan.VPCConfig.IsNull() {
		// no delete support
	} else if !plan.VPCConfig.IsUnknown() {
		stateVPCConfig := expandVPCConfig(ctx, state.VPCConfig, &resp.Diagnostics)
		planVPCConfig := expandVPCConfig(ctx, plan.VPCConfig, &resp.Diagnostics)

		if stateVPCConfig != nil && len(stateVPCConfig.SecurityGroupIds) > 0 &&
			planVPCConfig != nil && len(planVPCConfig.SecurityGroupIds) == 0 {
			attrsToDelete = append(attrsToDelete, awstypes.AppBlockBuilderAttributeVpcConfigurationSecurityGroupIds)
		}

		input.VpcConfig = planVPCConfig
	}

	input.AttributesToDelete = attrsToDelete

	if resp.Diagnostics.HasError() {
		return
	}

	restartAfter := false
	var err error
	mode := updateMode(state, plan)

	switch mode {
	case appBlockBuilderUpdateRequiresStop:
		restartAfter, err = r.ensureStopped(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream App Block Builder",
				fmt.Sprintf("Could not stop app block builder %q for update: %v", name, err),
			)
			return
		}
	case appBlockBuilderUpdateAllowedRunning:
		// proceed normally
	}

	out, err := r.appstreamClient.UpdateAppBlockBuilder(ctx, input)
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Updating AWS AppStream App Block Builder",
			fmt.Sprintf("Could not update app block builder %q: %v", name, err),
		)
		return
	}

	if out.AppBlockBuilder != nil && out.AppBlockBuilder.Arn != nil {
		_, tagDiags := r.tags.Apply(ctx, aws.ToString(out.AppBlockBuilder.Arn), plan.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if mode == appBlockBuilderUpdateRequiresStop && restartAfter {
		err = r.ensureRunning(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream App Block Builder",
				fmt.Sprintf("Could not start app block builder %q after successful update: %v", name, err),
			)
			return
		}
	}

	newState, diags := r.readAppBlockBuilder(ctx, plan)
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

func wasRunning(state awstypes.AppBlockBuilderState) bool {
	return state == awstypes.AppBlockBuilderStateRunning || state == awstypes.AppBlockBuilderStateStarting
}

func (r *resource) ensureStopped(ctx context.Context, name string) (restartAfter bool, err error) {
	stateCaptured := false

	err = util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			out, err := r.appstreamClient.DescribeAppBlockBuilders(ctx, &awsappstream.DescribeAppBlockBuildersInput{
				Names: []string{name},
			})
			if err != nil {
				if util.IsAppStreamNotFound(err) {
					return nil
				}
				return err
			}

			if len(out.AppBlockBuilders) == 0 {
				return nil
			}

			state := out.AppBlockBuilders[0].State

			if !stateCaptured {
				restartAfter = wasRunning(state)
				stateCaptured = true
			}

			switch state {
			case awstypes.AppBlockBuilderStateRunning:
				// stoppable state
				_, err = r.appstreamClient.StopAppBlockBuilder(ctx, &awsappstream.StopAppBlockBuilderInput{
					Name: aws.String(name),
				})
				if err != nil {
					if util.IsAppStreamNotFound(err) {
						return nil
					}
					return err
				}
				// retry as we just stopped the app block builder
				return fmt.Errorf("%w: current=%s", errUnexpectedAppBlockBuilderState, state)

			case awstypes.AppBlockBuilderStateStopped:
				// updatable state
				return nil

			default:
				return fmt.Errorf("%w: current=%s", errUnexpectedAppBlockBuilderState, state)
			}
		},
		util.WithTimeout(updateRetryTimeout),
		util.WithInitBackoff(updateRetryInitBackoff),
		util.WithMaxBackoff(updateRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeAppBlockBuilders.html
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_StopAppBlockBuilder.html
		util.WithRetryOnFns(
			isUnexpectedAppBlockBuilderStateError,
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
			out, err := r.appstreamClient.DescribeAppBlockBuilders(ctx, &awsappstream.DescribeAppBlockBuildersInput{
				Names: []string{name},
			})
			if err != nil {
				if util.IsAppStreamNotFound(err) {
					return nil
				}

				return err
			}

			if len(out.AppBlockBuilders) == 0 {
				return nil
			}

			state := out.AppBlockBuilders[0].State

			switch state {
			case awstypes.AppBlockBuilderStateStopped:
				// startable state
				_, err = r.appstreamClient.StartAppBlockBuilder(ctx, &awsappstream.StartAppBlockBuilderInput{
					Name: aws.String(name),
				})
				if err != nil {
					if util.IsAppStreamNotFound(err) {
						return nil
					}

					// app block builders cannot be started without an app block association.
					// if the association was removed externally, keep the builder stopped.
					if util.IsResourceNotAvailableException(err) &&
						strings.Contains(err.Error(), "Unassociated AppBlock Builder cannot be started.") {
						return nil
					}

					return err
				}
				// retry as we just started the app block builder
				return fmt.Errorf("%w: current=%s", errUnexpectedAppBlockBuilderState, state)

			case awstypes.AppBlockBuilderStateRunning:
				// desired state
				return nil

			default:
				return fmt.Errorf("%w: current=%s", errUnexpectedAppBlockBuilderState, state)
			}
		},
		util.WithTimeout(updateRetryTimeout),
		util.WithInitBackoff(updateRetryInitBackoff),
		util.WithMaxBackoff(updateRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeAppBlockBuilders.html
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_StartAppBlockBuilder.html
		util.WithRetryOnFns(
			isUnexpectedAppBlockBuilderStateError,
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
		),
	)
}
