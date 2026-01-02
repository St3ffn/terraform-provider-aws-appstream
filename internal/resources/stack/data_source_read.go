// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (ds *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config model

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
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
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

	state := &model{
		ID:                  types.StringValue(aws.ToString(stack.Name)),
		Name:                types.StringValue(aws.ToString(stack.Name)),
		Description:         util.StringOrNull(stack.Description),
		DisplayName:         util.StringOrNull(stack.DisplayName),
		StorageConnectors:   flattenStorageConnectorsData(ctx, stack.StorageConnectors, &resp.Diagnostics),
		RedirectURL:         util.StringOrNull(stack.RedirectURL),
		FeedbackURL:         util.StringOrNull(stack.FeedbackURL),
		UserSettings:        flattenUserSettingsData(ctx, stack.UserSettings, &resp.Diagnostics),
		ApplicationSettings: flattenApplicationSettingsData(ctx, stack.ApplicationSettings, &resp.Diagnostics),
		Tags:                types.Map{},
		AccessEndpoints:     flattenAccessEndpointsData(ctx, stack.AccessEndpoints, &resp.Diagnostics),
		EmbedHostDomains:    util.SetStringOrNull(ctx, stack.EmbedHostDomains, &resp.Diagnostics),
		StreamingExperienceSettings: flattenStreamingExperienceSettingsData(
			ctx, stack.StreamingExperienceSettings, &resp.Diagnostics,
		),
		ARN:         util.StringOrNull(stack.Arn),
		CreatedTime: util.StringFromTime(stack.CreatedTime),
		StackErrors: flattenStackErrorsData(ctx, stack.StackErrors, &resp.Diagnostics),
	}

	if !state.ARN.IsNull() {
		tags, diags := ds.tags.Read(ctx, state.ARN.ValueString())
		resp.Diagnostics.Append(diags...)
		state.Tags = tags
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
