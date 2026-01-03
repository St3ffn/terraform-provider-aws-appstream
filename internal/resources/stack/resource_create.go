// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan model

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

	input.Description = util.StringPointerOrNil(plan.Description)
	input.DisplayName = util.StringPointerOrNil(plan.DisplayName)
	input.RedirectURL = util.StringPointerOrNil(plan.RedirectURL)
	input.FeedbackURL = util.StringPointerOrNil(plan.FeedbackURL)

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

	input.EmbedHostDomains = util.ExpandStringSetOrNil(ctx, plan.EmbedHostDomains, &resp.Diagnostics)

	if !plan.StreamingExperienceSettings.IsNull() && !plan.StreamingExperienceSettings.IsUnknown() {
		input.StreamingExperienceSettings = expandStreamingExperienceSettings(
			ctx, plan.StreamingExperienceSettings, &resp.Diagnostics,
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var out *awsappstream.CreateStackOutput
	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			var err error
			out, err = r.appstreamClient.CreateStack(ctx, input)
			return err
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateStack.html
		util.WithRetryOnFns(
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
			util.IsResourceNotFoundException,
		),
	)

	if err != nil {
		if util.IsResourceAlreadyExists(err) {
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
