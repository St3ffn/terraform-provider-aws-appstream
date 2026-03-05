// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package fleet

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Delete calls the AWS delete/disassociate API and treats not-found
// responses as already deleted so destroy remains idempotent.
// Terraform state is then cleared by the framework lifecycle.
func (r *resource) Delete(ctx context.Context, req tfresource.DeleteRequest, resp *tfresource.DeleteResponse) {
	var state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if state.Name.IsNull() || state.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Cannot delete fleet because name must be known.",
		)
		return
	}

	name := state.Name.ValueString()

	_, err := r.ensureStopped(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream Fleet",
			fmt.Sprintf("Could not stop fleet %q for delete: %v", name, err),
		)
		return
	}

	_, err = r.appstreamClient.DeleteFleet(ctx, &awsappstream.DeleteFleetInput{
		Name: aws.String(name),
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
			"Error Deleting AWS AppStream Fleet",
			fmt.Sprintf("Could not delete fleet %q: %v", name, err),
		)
		return
	}
}
