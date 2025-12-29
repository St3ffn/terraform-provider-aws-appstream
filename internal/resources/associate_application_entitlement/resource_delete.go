// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_application_entitlement

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Delete(ctx context.Context, req tfresource.DeleteRequest, resp *tfresource.DeleteResponse) {
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	addDiagnostics(state, &resp.Diagnostics, diagnosticDelete)
	if resp.Diagnostics.HasError() {
		return
	}

	stackName := state.StackName.ValueString()
	entitlementName := state.EntitlementName.ValueString()
	applicationIdentifier := state.ApplicationIdentifier.ValueString()

	_, err := r.appstreamClient.DisassociateApplicationFromEntitlement(ctx, &awsappstream.DisassociateApplicationFromEntitlementInput{
		StackName:             aws.String(stackName),
		EntitlementName:       aws.String(entitlementName),
		ApplicationIdentifier: aws.String(applicationIdentifier),
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
			"Error Deleting AWS AppStream Application Entitlement Association",
			fmt.Sprintf("Could not disassociate application %q from entitlement %q (stack %q): %v",
				applicationIdentifier, entitlementName, stackName, err,
			),
		)
		return
	}
}
