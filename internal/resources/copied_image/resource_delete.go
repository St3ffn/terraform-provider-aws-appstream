// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package copied_image

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

	if state.DestinationImageName.IsNull() || state.DestinationImageName.IsUnknown() ||
		state.DestinationRegion.IsNull() || state.DestinationRegion.IsUnknown() ||
		state.SourceImageName.IsNull() || state.SourceImageName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Cannot delete copied image because destination_image_name, destination_region, and source_image_name must be known.",
		)
		return
	}

	name := state.DestinationImageName.ValueString()
	destinationRegion := state.DestinationRegion.ValueString()

	_, err := r.appstreamClient.DeleteImage(
		ctx,
		&awsappstream.DeleteImageInput{
			Name: aws.String(name),
		},
		func(o *awsappstream.Options) {
			o.Region = destinationRegion
		},
	)
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		// if it's already gone, that's fine for delete.
		if util.IsAppStreamNotFound(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream Copied Image",
			fmt.Sprintf("Could not delete copied image %q: %v", name, err),
		)
		return
	}
}
