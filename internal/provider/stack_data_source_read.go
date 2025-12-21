// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (ds *stackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config stackModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if config.Name.IsNull() || config.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Configuration",
			"Cannot read stack because name must be set and known.",
		)
		return
	}

	name := config.Name.ValueString()

	out, err := ds.appstreamClient.DescribeStacks(ctx, &awsappstream.DescribeStacksInput{
		Names: []string{name},
	})
	if err != nil {
		if isContextCanceled(ctx) {
			return
		}

		if isAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Stack Not Found",
				fmt.Sprintf("No stack named %q was found.", name),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream Stack",
			fmt.Sprintf("Could not read stack %q: %v", name, err),
		)
		return
	}

	if len(out.Stacks) == 0 {
		resp.Diagnostics.AddError(
			"AWS AppStream Stack Not Found",
			fmt.Sprintf("No stack named %q was found.", name),
		)
		return
	}

	stack := out.Stacks[0]
	if stack.Name == nil {
		resp.Diagnostics.AddError(
			"Unexpected AWS Response",
			fmt.Sprintf("Stack %q was returned without required identifiers.", name),
		)
		return
	}

	state := &stackModel{
		ID:                  types.StringValue(aws.ToString(stack.Name)),
		Name:                types.StringValue(aws.ToString(stack.Name)),
		Description:         stringOrNull(stack.Description),
		DisplayName:         stringOrNull(stack.DisplayName),
		StorageConnectors:   flattenStorageConnectors(ctx, stack.StorageConnectors, &resp.Diagnostics),
		RedirectURL:         stringOrNull(stack.RedirectURL),
		FeedbackURL:         stringOrNull(stack.FeedbackURL),
		UserSettings:        flattenUserSettings(ctx, stack.UserSettings, &resp.Diagnostics),
		ApplicationSettings: flattenApplicationSettings(ctx, stack.ApplicationSettings, &resp.Diagnostics),
		Tags:                types.Map{},
		AccessEndpoints:     flattenAccessEndpoints(ctx, stack.AccessEndpoints, &resp.Diagnostics),
		EmbedHostDomains:    setStringOrNull(ctx, stack.EmbedHostDomains, &resp.Diagnostics),
		StreamingExperienceSettings: flattenStreamingExperienceSettings(
			ctx, stack.StreamingExperienceSettings, &resp.Diagnostics,
		),
		ARN:         stringOrNull(stack.Arn),
		CreatedTime: stringFromTime(stack.CreatedTime),
		StackErrors: flattenStackErrors(ctx, stack.StackErrors, &resp.Diagnostics),
	}

	if !state.ARN.IsNull() {
		tags, diags := ds.tags.Read(ctx, state.ARN.ValueString())
		resp.Diagnostics.Append(diags...)
		state.Tags = tags
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
