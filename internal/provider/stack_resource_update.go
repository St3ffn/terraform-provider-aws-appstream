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

func (r *stackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan stackModel
	var state stackModel

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

	optionalStringUpdate(plan.Description, state.Description, func(description *string) {
		input.Description = description
	})
	optionalStringUpdate(plan.DisplayName, state.DisplayName, func(displayName *string) {
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
			input.EmbedHostDomains = expandStringSetOrNil(
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
		if isContextCanceled(ctx) {
			return
		}

		if isAppStreamNotFound(err) {
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

	newState, diags := r.readStack(ctx, name)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if isContextCanceled(ctx) {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
