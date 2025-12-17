// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *associateApplicationEntitlementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state associateApplicationEntitlementModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	addAssocPartsDiagnostics(state, &resp.Diagnostics, assocDiagDelete)
	if resp.Diagnostics.HasError() {
		return
	}

	stackName := state.StackName.ValueString()
	entName := state.EntitlementName.ValueString()
	appID := state.ApplicationIdentifier.ValueString()

	_, err := r.appstreamClient.DisassociateApplicationFromEntitlement(ctx, &awsappstream.DisassociateApplicationFromEntitlementInput{
		StackName:             aws.String(stackName),
		EntitlementName:       aws.String(entName),
		ApplicationIdentifier: aws.String(appID),
	})
	if err != nil {
		// respect cancellation/deadlines
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}

		// if it's already gone, that's fine for Delete.
		if isAppStreamNotFound(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream Application Entitlement Association",
			fmt.Sprintf("Could not disassociate application %q from entitlement %q (stack %q): %v",
				appID, entName, stackName, err,
			),
		)
		return
	}
}
