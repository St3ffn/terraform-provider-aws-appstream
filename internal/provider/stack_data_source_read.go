// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
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
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
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

	var state stackModel

	state.Name = types.StringValue(aws.ToString(stack.Name))
	state.ID = types.StringValue(aws.ToString(stack.Name))

	state.Description = stringOrNull(stack.Description)
	state.DisplayName = stringOrNull(stack.DisplayName)
	state.RedirectURL = stringOrNull(stack.RedirectURL)
	state.FeedbackURL = stringOrNull(stack.FeedbackURL)
	state.StorageConnectors = flattenStorageConnectors(ctx, stack.StorageConnectors, &resp.Diagnostics)
	state.UserSettings = flattenUserSettings(ctx, stack.UserSettings, &resp.Diagnostics)
	state.ApplicationSettings = flattenApplicationSettings(ctx, stack.ApplicationSettings, &resp.Diagnostics)
	state.AccessEndpoints = flattenAccessEndpoints(ctx, stack.AccessEndpoints, &resp.Diagnostics)
	state.EmbedHostDomains = setStringOrNull(ctx, stack.EmbedHostDomains, &resp.Diagnostics)
	state.StreamingExperienceSettings = flattenStreamingExperienceSettings(ctx, stack.StreamingExperienceSettings, &resp.Diagnostics)

	state.ARN = stringOrNull(stack.Arn)
	state.CreatedTime = stringFromTime(stack.CreatedTime)
	state.StackErrors = flattenStackErrors(ctx, stack.StackErrors, &resp.Diagnostics)

	if !state.ARN.IsNull() {
		tags, diags := readTags(ctx, ds.taggingClient, state.ARN.ValueString())
		resp.Diagnostics.Append(diags...)
		state.Tags = flattenTags(ctx, tags, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
