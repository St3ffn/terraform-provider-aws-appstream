// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package application

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

	if plan.ID.IsNull() || plan.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update application because id must be known.",
		)
		return
	}

	diff := newResourceDiff(state, plan)
	arn := plan.ID.ValueString()

	name, err := applicationNameFromARN(arn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			fmt.Sprintf("Could not parse application name from ARN %q: %v", arn, err),
		)
		return
	}

	if diff.HasRemoteChanges() {
		input := &awsappstream.UpdateApplicationInput{
			Name: aws.String(name),
		}

		var attrsToDelete []awstypes.ApplicationAttribute

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
		if diff.LaunchPath.IsChanged() {
			util.OptionalStringUpdate(plan.LaunchPath, state.LaunchPath, func(v *string) {
				input.LaunchPath = v
			})
		}

		if diff.AppBlockARN.IsChanged() {
			util.OptionalStringUpdate(plan.AppBlockARN, state.AppBlockARN, func(v *string) {
				input.AppBlockArn = v
			})
		}

		if diff.WorkingDirectory.IsChanged() {
			if plan.WorkingDirectory.IsNull() {
				attrsToDelete = append(attrsToDelete, awstypes.ApplicationAttributeWorkingDirectory)
			} else {
				input.WorkingDirectory = plan.WorkingDirectory.ValueStringPointer()
			}
		}

		if diff.LaunchParameters.IsChanged() {
			if plan.LaunchParameters.IsNull() {
				attrsToDelete = append(attrsToDelete, awstypes.ApplicationAttributeLaunchParameters)
			} else {
				input.LaunchParameters = plan.LaunchParameters.ValueStringPointer()
			}
		}

		if diff.IconS3Location.IsChanged() {
			if plan.IconS3Location.IsNull() {
				// no delete support
			} else {
				input.IconS3Location = expandIconS3Location(
					ctx, plan.IconS3Location, &resp.Diagnostics,
				)
			}
		}

		if len(attrsToDelete) > 0 {
			input.AttributesToDelete = attrsToDelete
		}

		if resp.Diagnostics.HasError() {
			return
		}

		_, err = r.appstreamClient.UpdateApplication(ctx, input)
		if err != nil {
			if util.IsContextCanceled(err) {
				return
			}

			if util.IsAppStreamNotFound(err) {
				resp.State.RemoveResource(ctx)
				return
			}

			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Application",
				fmt.Sprintf("Could not update application %q: %v", arn, err),
			)
			return
		}
	}

	if diff.RequiresTagApply() {
		_, tagDiags := r.tags.Apply(ctx, arn, plan.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	newState, diags := r.readApplication(ctx, plan)
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
