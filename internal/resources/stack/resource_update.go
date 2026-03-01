// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package stack

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
			"Cannot update stack because name must be known.",
		)
		return
	}

	diff := newResourceDiff(state, plan)
	name := plan.Name.ValueString()

	arnForTags := ""
	switch {
	case !plan.ARN.IsNull() && !plan.ARN.IsUnknown():
		arnForTags = plan.ARN.ValueString()
	case !state.ARN.IsNull() && !state.ARN.IsUnknown():
		arnForTags = state.ARN.ValueString()
	}

	if diff.HasRemoteChanges() {
		input := &awsappstream.UpdateStackInput{
			Name: aws.String(name),
		}
		var attrToDelete []awstypes.StackAttribute

		if diff.Description.IsChanged() {
			util.OptionalStringUpdate(plan.Description, state.Description, func(description *string) {
				input.Description = description
			})
		}
		if diff.DisplayName.IsChanged() {
			util.OptionalStringUpdate(plan.DisplayName, state.DisplayName, func(displayName *string) {
				input.DisplayName = displayName
			})
		}

		if diff.RedirectUrl.IsChanged() {
			if plan.RedirectUrl.IsNull() {
				attrToDelete = append(attrToDelete, awstypes.StackAttributeRedirectUrl)
			} else {
				input.RedirectURL = aws.String(plan.RedirectUrl.ValueString())
			}
		}

		if diff.FeedbackUrl.IsChanged() {
			if plan.FeedbackUrl.IsNull() {
				attrToDelete = append(attrToDelete, awstypes.StackAttributeFeedbackUrl)
			} else {
				input.FeedbackURL = aws.String(plan.FeedbackUrl.ValueString())
			}
		}

		if diff.StorageConnectors.IsChanged() {
			if plan.StorageConnectors.IsNull() {
				attrToDelete = append(attrToDelete, awstypes.StackAttributeStorageConnectors)
			} else {
				attrToDelete = append(
					attrToDelete,
					storageConnectorAttributesToDelete(ctx, state.StorageConnectors, plan.StorageConnectors, &resp.Diagnostics)...,
				)

				input.StorageConnectors = expandStorageConnectors(
					ctx,
					plan.StorageConnectors,
					&resp.Diagnostics,
				)
			}
		}

		if diff.UserSettings.IsChanged() {
			if plan.UserSettings.IsNull() {
				attrToDelete = append(attrToDelete, awstypes.StackAttributeUserSettings)
			} else {
				input.UserSettings = expandUserSettings(
					ctx,
					plan.UserSettings,
					&resp.Diagnostics,
				)
			}
		}

		if diff.ApplicationSettings.IsChanged() {
			if plan.ApplicationSettings.IsNull() {
				input.ApplicationSettings = nil
			} else {
				input.ApplicationSettings = expandApplicationSettings(
					ctx,
					plan.ApplicationSettings,
					&resp.Diagnostics,
				)
			}
		}

		if diff.AccessEndpoints.IsChanged() {
			if plan.AccessEndpoints.IsNull() {
				attrToDelete = append(attrToDelete, awstypes.StackAttributeAccessEndpoints)
			} else {
				input.AccessEndpoints = expandAccessEndpoints(
					ctx,
					plan.AccessEndpoints,
					&resp.Diagnostics,
				)
			}
		}

		if diff.EmbedHostDomains.IsChanged() {
			if plan.EmbedHostDomains.IsNull() {
				attrToDelete = append(attrToDelete, awstypes.StackAttributeEmbedHostDomains)
			} else {
				input.EmbedHostDomains = util.ExpandStringSetOrNil(
					ctx,
					plan.EmbedHostDomains,
					&resp.Diagnostics,
				)
			}
		}

		if diff.StreamingExperienceSettings.IsChanged() {
			if plan.StreamingExperienceSettings.IsNull() {
				attrToDelete = append(attrToDelete, awstypes.StackAttributeStreamingExperienceSettings)
			} else {
				input.StreamingExperienceSettings = expandStreamingExperienceSettings(
					ctx,
					plan.StreamingExperienceSettings,
					&resp.Diagnostics,
				)
			}
		}

		input.AttributesToDelete = attrToDelete

		if resp.Diagnostics.HasError() {
			return
		}

		out, err := r.appstreamClient.UpdateStack(ctx, input)
		if err != nil {
			if util.IsContextCanceled(err) {
				return
			}

			if util.IsAppStreamNotFound(err) {
				resp.State.RemoveResource(ctx)
				return
			}

			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Stack",
				fmt.Sprintf("Could not update stack %q: %v", name, err),
			)
			return
		}

		if out.Stack != nil && out.Stack.Arn != nil {
			arnForTags = aws.ToString(out.Stack.Arn)
		}
	}

	if diff.RequiresTagApply() {
		if arnForTags == "" {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Stack",
				fmt.Sprintf("Could not apply tags for stack %q: missing ARN in state and plan", name),
			)
			return
		}
		_, tagDiags := r.tags.Apply(ctx, arnForTags, plan.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	newState, diags := r.readStack(ctx, plan)
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
