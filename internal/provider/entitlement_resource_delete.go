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

func (r *entitlementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entitlementModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if state.StackName.IsNull() || state.StackName.IsUnknown() ||
		state.Name.IsNull() || state.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Cannot delete entitlement because stack_name, and name must be known.",
		)
		return
	}

	stackName := state.StackName.ValueString()
	name := state.Name.ValueString()

	_, err := r.appstreamClient.DeleteEntitlement(ctx, &awsappstream.DeleteEntitlementInput{
		StackName: aws.String(stackName),
		Name:      aws.String(name),
	})
	if err != nil {
		if isContextCanceled(ctx) {
			return
		}

		// if it's already gone, that's fine for Delete.
		if isAppStreamNotFound(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream Entitlement",
			fmt.Sprintf("Could not delete entitlement %q in stack %q: %v", name, stackName, err),
		)
		return
	}
}
