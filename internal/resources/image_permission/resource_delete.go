// Copyright St3ffn 2025, 2026
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

func (r *resource) Delete(ctx context.Context, req tfresource.DeleteRequest, resp *tfresource.DeleteResponse) {
	var state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if state.Name.IsNull() || state.Name.IsUnknown() ||
		state.SharedAccountID.IsNull() || state.SharedAccountID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Cannot delete image permission because name, and shared_account_id must be known.",
		)
		return
	}

	name := state.Name.ValueString()
	sharedAccountID := state.SharedAccountID.ValueString()

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.DeleteImagePermissions(ctx, &awsappstream.DeleteImagePermissionsInput{
				Name:            aws.String(name),
				SharedAccountId: aws.String(sharedAccountID),
			})
			return err
		},
		util.WithTimeout(retryTimeout),
		util.WithInitBackoff(retryInitBackoff),
		util.WithMaxBackoff(retryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DeleteImagePermissions.html
		util.WithRetryOnFns(
			util.IsResourceNotAvailableException,
		),
	)

	if err != nil {
		// if it's already gone, that's fine for delete.
		if util.IsAppStreamNotFound(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream Image Permission",
			fmt.Sprintf("Could not delete image permission for image %q shared with account %q: %v", name, sharedAccountID, err),
		)
		return
	}
}
