// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *stackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state stackModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if state.Name.IsNull() || state.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Required attributes name is missing from state. "+
				"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
		)
		return
	}

	name := state.Name.ValueString()

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

func (r *stackResource) readStack(ctx context.Context, name string) (*stackModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	out, err := r.appstreamClient.DescribeStacks(ctx, &awsappstream.DescribeStacksInput{
		Names: []string{name},
	})
	if err != nil {
		if isContextCanceled(ctx) {
			return nil, diags
		}

		if isAppStreamNotFound(err) {
			return nil, diags
		}
		diags.AddError(
			"Error Reading AWS AppStream Stack",
			fmt.Sprintf("Could not read stack %q: %v", name, err),
		)
		return nil, diags
	}

	if len(out.Stacks) == 0 {
		return nil, diags
	}

	stack := out.Stacks[0]
	if stack.Name == nil {
		return nil, diags
	}

	state := &stackModel{
		ID:                          types.StringValue(aws.ToString(stack.Name)),
		Name:                        types.StringValue(aws.ToString(stack.Name)),
		Description:                 stringOrNull(stack.Description),
		DisplayName:                 stringOrNull(stack.DisplayName),
		StorageConnectors:           flattenStorageConnectors(ctx, stack.StorageConnectors, &diags),
		RedirectURL:                 stringOrNull(stack.RedirectURL),
		FeedbackURL:                 stringOrNull(stack.FeedbackURL),
		UserSettings:                flattenUserSettings(ctx, stack.UserSettings, &diags),
		ApplicationSettings:         flattenApplicationSettings(ctx, stack.ApplicationSettings, &diags),
		Tags:                        types.Map{},
		AccessEndpoints:             flattenAccessEndpoints(ctx, stack.AccessEndpoints, &diags),
		EmbedHostDomains:            setStringOrNull(ctx, stack.EmbedHostDomains, &diags),
		StreamingExperienceSettings: flattenStreamingExperienceSettings(ctx, stack.StreamingExperienceSettings, &diags),
		ARN:                         stringOrNull(stack.Arn),
		CreatedTime:                 stringFromTime(stack.CreatedTime),
		StackErrors:                 flattenStackErrors(ctx, stack.StackErrors, &diags),
	}

	if !state.ARN.IsNull() {
		tags, tagDiags := r.tags.Read(ctx, state.ARN.ValueString())
		diags.Append(tagDiags...)
		state.Tags = tags
	}

	if diags.HasError() {
		return nil, diags
	}
	return state, diags
}
