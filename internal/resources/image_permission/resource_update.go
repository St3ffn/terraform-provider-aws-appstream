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

	if plan.Name.IsNull() || plan.Name.IsUnknown() ||
		plan.SharedAccountID.IsNull() || plan.SharedAccountID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update image permission because name and shared_account_id must be known.",
		)
		return
	}

	name := plan.Name.ValueString()
	sharedAccountID := plan.SharedAccountID.ValueString()

	// guard against unexpected identity drift
	if !state.Name.IsNull() && !state.Name.IsUnknown() &&
		!state.SharedAccountID.IsNull() && !state.SharedAccountID.IsUnknown() {
		if state.Name.ValueString() != name || state.SharedAccountID.ValueString() != sharedAccountID {
			resp.Diagnostics.AddError(
				"Unexpected Update Request",
				"Image permission identity (name|shared_account_id) changed during update. This should trigger replacement. Please report this issue.",
			)
			return
		}
	}

	input := &awsappstream.UpdateImagePermissionsInput{
		Name:             aws.String(name),
		SharedAccountId:  aws.String(sharedAccountID),
		ImagePermissions: expandImagePermissions(ctx, plan.ImagePermissions, &resp.Diagnostics),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.UpdateImagePermissions(ctx, input)
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
		// disappeared, treat as gone and next plan/apply will recreate
		if util.IsAppStreamNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Updating AWS AppStream Image Permission",
			fmt.Sprintf("Could not update image permission for image %q shared with account %q: %v", name, sharedAccountID, err),
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
