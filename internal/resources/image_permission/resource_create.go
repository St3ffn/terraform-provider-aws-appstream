// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_permission

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	if plan.Name.IsNull() || plan.Name.IsUnknown() ||
		plan.SharedAccountID.IsNull() || plan.SharedAccountID.IsUnknown() ||
		plan.ImagePermissions.IsNull() || plan.ImagePermissions.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create image permission because name, shared_account_id, and image_permissions must be known.",
		)
		return
	}

	name := plan.Name.ValueString()
	sharedAccountID := plan.SharedAccountID.ValueString()

	input := &awsappstream.UpdateImagePermissionsInput{
		Name:            aws.String(name),
		SharedAccountId: aws.String(sharedAccountID),
	}

	if !plan.ImagePermissions.IsNull() && !plan.ImagePermissions.IsUnknown() {
		input.ImagePermissions = expandImagePermissions(ctx, plan.ImagePermissions, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			var err error
			_, err = r.appstreamClient.UpdateImagePermissions(ctx, input)
			return err
		},
		util.WithTimeout(retryTimeout),
		util.WithInitBackoff(retryInitBackoff),
		util.WithMaxBackoff(retryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_UpdateImagePermissions.html
		util.WithRetryOnFns(
			util.IsResourceNotAvailableException,
		),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Image Permission",
			fmt.Sprintf("Could not create image permission for image %q shared with account %q: %v", name, sharedAccountID, err),
		)
		return
	}

	newState, diags := r.readImagePermission(ctx, plan)
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
