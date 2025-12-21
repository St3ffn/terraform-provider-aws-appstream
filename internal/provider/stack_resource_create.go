// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *stackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan stackModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create stack because name must be known.",
		)
		return
	}

	name := plan.Name.ValueString()

	input := &awsappstream.CreateStackInput{
		Name: aws.String(name),
	}

	input.Description = stringPointerOrNil(plan.Description)
	input.DisplayName = stringPointerOrNil(plan.DisplayName)
	input.RedirectURL = stringPointerOrNil(plan.RedirectURL)
	input.FeedbackURL = stringPointerOrNil(plan.FeedbackURL)

	if !plan.StorageConnectors.IsNull() && !plan.StorageConnectors.IsUnknown() {
		input.StorageConnectors = expandStorageConnectors(ctx, plan.StorageConnectors, &resp.Diagnostics)
	}

	if !plan.UserSettings.IsNull() && !plan.UserSettings.IsUnknown() {
		input.UserSettings = expandUserSettings(ctx, plan.UserSettings, &resp.Diagnostics)
	}

	if !plan.ApplicationSettings.IsNull() && !plan.ApplicationSettings.IsUnknown() {
		input.ApplicationSettings = expandApplicationSettings(ctx, plan.ApplicationSettings, &resp.Diagnostics)
	}

	if !plan.AccessEndpoints.IsNull() && !plan.AccessEndpoints.IsUnknown() {
		input.AccessEndpoints = expandAccessEndpoints(ctx, plan.AccessEndpoints, &resp.Diagnostics)
	}

	input.EmbedHostDomains = expandStringSetOrNil(ctx, plan.EmbedHostDomains, &resp.Diagnostics)

	if !plan.StreamingExperienceSettings.IsNull() && !plan.StreamingExperienceSettings.IsUnknown() {
		input.StreamingExperienceSettings = expandStreamingExperienceSettings(
			ctx, plan.StreamingExperienceSettings, &resp.Diagnostics,
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.appstreamClient.CreateStack(ctx, input)
	if err != nil {
		if isContextCanceled(ctx) {
			return
		}

		if isAppStreamAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Stack Already Exists",
				fmt.Sprintf(
					"A stack named %q already exists. To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> %q",
					name, name,
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Stack",
			fmt.Sprintf("Could not create stack %q: %v", name, err),
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
