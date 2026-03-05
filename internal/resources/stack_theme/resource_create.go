// Copyright St3ffn 2025, 2026
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

// Create validates plan input, performs the required AWS create/associate calls,
// retries configured transient API errors, and then reads the remote object back
// so state is written from the authoritative AWS response.
func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	if plan.StackName.IsNull() || plan.StackName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create stack theme because stack_name must be known.",
		)
		return
	}

	stackName := plan.StackName.ValueString()

	input := &awsappstream.CreateThemeForStackInput{
		StackName:                  aws.String(stackName),
		TitleText:                  util.StringPointerOrNil(plan.TitleText),
		ThemeStyling:               awstypes.ThemeStyling(plan.ThemeStyling.ValueString()),
		OrganizationLogoS3Location: expandS3Location(ctx, plan.OrganizationLogoS3Location, &resp.Diagnostics),
		FaviconS3Location:          expandS3Location(ctx, plan.FaviconS3Location, &resp.Diagnostics),
	}

	if !plan.FooterLinks.IsNull() && !plan.FooterLinks.IsUnknown() {
		input.FooterLinks = expandFooterLinks(ctx, plan.FooterLinks, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.CreateThemeForStack(ctx, input)
			return err
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateThemeForStack.html
		util.WithRetryOnFns(
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
			util.IsResourceNotFoundException,
		),
	)

	if err != nil {
		if util.IsResourceAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Stack Theme Already Exists",
				fmt.Sprintf(
					"A theme for stack %q already exists. To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> %q",
					stackName, stackName,
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Stack Theme",
			fmt.Sprintf("Could not create theme for stack %q: %v", stackName, err),
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
