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

func (r *associateFleetStackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state associateFleetStackModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	addAssociateFleetStackDiagnostics(state, &resp.Diagnostics, associateDiagnosticDelete)
	if resp.Diagnostics.HasError() {
		return
	}

	fleetName := state.FleetName.ValueString()
	stackName := state.StackName.ValueString()

	_, err := r.appstreamClient.DisassociateFleet(ctx, &awsappstream.DisassociateFleetInput{
		FleetName: aws.String(fleetName),
		StackName: aws.String(stackName),
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
			"Error Deleting AWS AppStream Fleet Stack Association",
			fmt.Sprintf("Could not disassociate fleet %q from stack %q: %v", fleetName, stackName, err),
		)
		return
	}
}
