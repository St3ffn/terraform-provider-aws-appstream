// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Read(ctx context.Context, req tfresource.ReadRequest, resp *tfresource.ReadResponse) {
	var state model

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

	newState, diags := r.readStack(ctx, state)
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

func (r *resource) readStack(ctx context.Context, prior model) (*model, diag.Diagnostics) {
	var diags diag.Diagnostics

	name := prior.Name.ValueString()

	out, err := r.appstreamClient.DescribeStacks(ctx, &awsappstream.DescribeStacksInput{
		Names: []string{name},
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		if util.IsAppStreamNotFound(err) {
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

	state := &model{
		ID:                          types.StringValue(aws.ToString(stack.Name)),
		Name:                        types.StringValue(aws.ToString(stack.Name)),
		Description:                 util.StringOrNull(stack.Description),
		DisplayName:                 util.StringOrNull(stack.DisplayName),
		StorageConnectors:           flattenStorageConnectorsResource(ctx, prior.StorageConnectors, stack.StorageConnectors, &diags),
		RedirectURL:                 util.StringOrNull(stack.RedirectURL),
		FeedbackURL:                 util.StringOrNull(stack.FeedbackURL),
		UserSettings:                flattenUserSettingsResource(ctx, prior.UserSettings, stack.UserSettings, &diags),
		ApplicationSettings:         flattenApplicationSettingsResource(ctx, prior.ApplicationSettings, stack.ApplicationSettings, &diags),
		Tags:                        types.Map{},
		AccessEndpoints:             flattenAccessEndpointsResource(ctx, prior.AccessEndpoints, stack.AccessEndpoints, &diags),
		EmbedHostDomains:            util.FlattenOwnedStringSet(ctx, prior.EmbedHostDomains, stack.EmbedHostDomains, &diags),
		StreamingExperienceSettings: flattenStreamingExperienceSettingsResource(ctx, prior.StreamingExperienceSettings, stack.StreamingExperienceSettings, &diags),
		ARN:                         util.StringOrNull(stack.Arn),
		CreatedTime:                 util.StringFromTime(stack.CreatedTime),
		StackErrors:                 flattenStackErrorsData(ctx, stack.StackErrors, &diags),
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
