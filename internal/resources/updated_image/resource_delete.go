// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package updated_image

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Delete calls the AWS delete API and treats not-found responses as already deleted
// so destroy remains idempotent. Terraform state is then cleared by the framework lifecycle.
func (r *resource) Delete(ctx context.Context, req tfresource.DeleteRequest, resp *tfresource.DeleteResponse) {
	var state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if state.ExistingImageName.IsNull() || state.ExistingImageName.IsUnknown() ||
		state.NewImageName.IsNull() || state.NewImageName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Cannot delete updated image because existing_image_name and new_image_name must be known.",
		)
		return
	}

	newImageName := state.NewImageName.ValueString()

	_, err := r.appstreamClient.DeleteImage(ctx, &awsappstream.DeleteImageInput{
		Name: aws.String(newImageName),
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		// if it's already gone, that's fine for delete.
		if util.IsAppStreamNotFound(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream Updated Image",
			fmt.Sprintf("Could not delete updated image %q: %v", newImageName, err),
		)
		return
	}
}
