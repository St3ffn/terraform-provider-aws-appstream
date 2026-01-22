// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack_theme

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
	var config dataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if config.StackName.IsNull() || config.StackName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Configuration",
			"Cannot read stack theme because stack_name must be set and known.",
		)
		return
	}

	stackName := config.StackName.ValueString()

	out, err := ds.appstreamClient.DescribeThemeForStack(ctx, &awsappstream.DescribeThemeForStackInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Stack Theme Not Found",
				fmt.Sprintf("No theme for stack %q was found.", stackName),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream Stack Theme",
			fmt.Sprintf("Could not read theme for stack %q: %v", stackName, err),
		)
		return
	}

	if out.Theme == nil {
		resp.Diagnostics.AddError(
			"AWS AppStream Stack Theme Not Found",
			fmt.Sprintf("No theme for stack %q was found.", stackName),
		)
		return
	}

	theme := out.Theme
	if theme.StackName == nil {
		resp.Diagnostics.AddError(
			"Unexpected AWS Response",
			fmt.Sprintf("Theme for stack %q was returned without required identifiers.", stackName),
		)
		return
	}

	state := &dataSourceModel{
		ID:                       types.StringValue(aws.ToString(theme.StackName)),
		StackName:                types.StringValue(aws.ToString(theme.StackName)),
		TitleText:                util.StringOrNull(theme.ThemeTitleText),
		ThemeStyling:             types.StringValue(string(theme.ThemeStyling)),
		FooterLinks:              flattenFooterLinks(ctx, theme.ThemeFooterLinks, &resp.Diagnostics),
		State:                    types.StringValue(string(theme.State)),
		CreatedTime:              util.StringFromTime(theme.CreatedTime),
		ThemeOrganizationLogoURL: util.StringOrNull(theme.ThemeOrganizationLogoURL),
		ThemeFaviconURL:          util.StringOrNull(theme.ThemeFaviconURL),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
