// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *stackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state stackModel

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
			"Cannot delete stack because name must be known.",
		)
		return
	}

	name := state.Name.ValueString()

	_, err := r.appstreamClient.DeleteStack(ctx, &awsappstream.DeleteStackInput{
		Name: aws.String(name),
	})
	if err != nil {
		if isContextCanceled(ctx) {
			return
		}

		// if it's already gone, that's fine for delete.
		if isAppStreamNotFound(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream Stack",
			fmt.Sprintf("Could not delete stack %q: %v", name, err),
		)
		return
	}
}
