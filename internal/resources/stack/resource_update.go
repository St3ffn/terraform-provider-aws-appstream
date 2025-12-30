// Copyright (c) St3ffn
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
	var plan model
	var state model

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

	name := plan.Name.ValueString()

	// guard against unexpected identity drift
	if !state.Name.IsNull() && !state.Name.IsUnknown() {
		if state.Name.ValueString() != name {
			resp.Diagnostics.AddError(
				"Unexpected Update Request",
				"Stack identity (name) changed during update. This should trigger replacement. Please report this issue.",
			)
			return
		}
	}

	input := &awsappstream.UpdateStackInput{
		Name: aws.String(name),
	}
	var attrToDelete []awstypes.StackAttribute

	util.OptionalStringUpdate(plan.Description, state.Description, func(description *string) {
		input.Description = description
	})
	util.OptionalStringUpdate(plan.DisplayName, state.DisplayName, func(displayName *string) {
		input.DisplayName = displayName
	})

	if !plan.RedirectURL.IsUnknown() {
		if plan.RedirectURL.IsNull() {
			attrToDelete = append(attrToDelete, awstypes.StackAttributeRedirectUrl)
		} else {
			input.RedirectURL = aws.String(plan.RedirectURL.ValueString())
		}
	}

	if !plan.FeedbackURL.IsUnknown() {
		if plan.FeedbackURL.IsNull() {
			attrToDelete = append(attrToDelete, awstypes.StackAttributeFeedbackUrl)
		} else {
			input.FeedbackURL = aws.String(plan.FeedbackURL.ValueString())
		}
	}

	if !plan.StorageConnectors.IsUnknown() {
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

	if !plan.UserSettings.IsUnknown() {
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

	if !plan.ApplicationSettings.IsUnknown() {
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

	if !plan.AccessEndpoints.IsUnknown() {
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

	if !plan.EmbedHostDomains.IsUnknown() {
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

	if !plan.StreamingExperienceSettings.IsUnknown() {
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
		_, tagDiags := r.tags.Apply(ctx, aws.ToString(out.Stack.Arn), plan.Tags)
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
		if util.IsContextCanceled(ctx.Err()) {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
