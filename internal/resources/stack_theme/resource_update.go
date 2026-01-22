// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack_theme

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Update(ctx context.Context, req tfresource.UpdateRequest, resp *tfresource.UpdateResponse) {
	var plan resourceModel
	var state resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if plan.StackName.IsNull() || plan.StackName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update stack theme because stack_name must be known.",
		)
		return
	}

	stackName := plan.StackName.ValueString()

	input := &awsappstream.UpdateThemeForStackInput{
		StackName: aws.String(stackName),
	}
	var attrToDelete []awstypes.ThemeAttribute

	util.OptionalStringUpdate(plan.TitleText, state.TitleText, func(titleText *string) {
		input.TitleText = titleText
	})

	if plan.ThemeStyling.IsNull() {
		// no delete support
	} else if !plan.ThemeStyling.IsUnknown() {
		input.ThemeStyling = awstypes.ThemeStyling(plan.ThemeStyling.ValueString())
	}

	if plan.OrganizationLogoS3Location.IsNull() {
		// no delete support
	} else if !plan.OrganizationLogoS3Location.IsUnknown() {
		input.OrganizationLogoS3Location = expandS3Location(ctx, plan.OrganizationLogoS3Location, &resp.Diagnostics)
	}

	if plan.FaviconS3Location.IsNull() {
		// no delete support
	} else if !plan.FaviconS3Location.IsUnknown() {
		input.FaviconS3Location = expandS3Location(ctx, plan.FaviconS3Location, &resp.Diagnostics)
	}

	if !plan.FooterLinks.IsUnknown() {
		if plan.FooterLinks.IsNull() {
			attrToDelete = append(attrToDelete, awstypes.ThemeAttributeFooterLinks)
		} else {
			input.FooterLinks = expandFooterLinks(ctx, plan.FooterLinks, &resp.Diagnostics)
		}
	}

	input.AttributesToDelete = attrToDelete

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.appstreamClient.UpdateThemeForStack(ctx, input)
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Updating AWS AppStream Stack Theme",
			fmt.Sprintf("Could not update theme for stack %q: %v", stackName, err),
		)
		return
	}

	newState, diags := r.readStackTheme(ctx, plan)
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
