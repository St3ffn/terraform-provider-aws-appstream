// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_theme

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
	var state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if state.StackName.IsNull() || state.StackName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Required attributes stack_name is missing from state. "+
				"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
		)
		return
	}

	newState, diags := r.readStackTheme(ctx, state)
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

func (r *resource) readStackTheme(ctx context.Context, prior resourceModel) (*resourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	stackName := prior.StackName.ValueString()

	out, err := r.appstreamClient.DescribeThemeForStack(ctx, &awsappstream.DescribeThemeForStackInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		if util.IsAppStreamNotFound(err) {
			return nil, diags
		}
		diags.AddError(
			"Error Reading AWS AppStream Stack Theme",
			fmt.Sprintf("Could not read theme for stack %q: %v", stackName, err),
		)
		return nil, diags
	}

	if out.Theme == nil {
		return nil, diags
	}

	theme := out.Theme
	if theme.StackName == nil {
		return nil, diags
	}

	state := &resourceModel{
		ID:                         types.StringValue(aws.ToString(theme.StackName)),
		StackName:                  types.StringValue(aws.ToString(theme.StackName)),
		TitleText:                  util.StringOrNull(theme.ThemeTitleText),
		ThemeStyling:               types.StringValue(string(theme.ThemeStyling)),
		OrganizationLogoS3Location: prior.OrganizationLogoS3Location,
		FaviconS3Location:          prior.FaviconS3Location,
		FooterLinks:                flattenFooterLinks(ctx, theme.ThemeFooterLinks, &diags),
		State:                      types.StringValue(string(theme.State)),
		CreatedTime:                util.StringFromTime(theme.CreatedTime),
		ThemeOrganizationLogoUrl:   util.StringOrNull(theme.ThemeOrganizationLogoURL),
		ThemeFaviconUrl:            util.StringOrNull(theme.ThemeFaviconURL),
	}

	if diags.HasError() {
		return nil, diags
	}
	return state, diags
}
