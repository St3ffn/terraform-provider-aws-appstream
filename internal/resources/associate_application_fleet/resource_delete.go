// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_application_fleet

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

	addDiagnostics(state, &resp.Diagnostics, diagnosticDelete)
	if resp.Diagnostics.HasError() {
		return
	}

	fleetName := state.FleetName.ValueString()
	applicationARN := state.ApplicationARN.ValueString()

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.DisassociateApplicationFleet(
				ctx,
				&awsappstream.DisassociateApplicationFleetInput{
					FleetName:      aws.String(fleetName),
					ApplicationArn: aws.String(applicationARN),
				},
			)
			return err
		},
		util.WithTimeout(deleteRetryTimeout),
		util.WithInitBackoff(deleteRetryInitBackoff),
		util.WithMaxBackoff(deleteRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DisassociateApplicationFleet.html
		util.WithRetryOnFns(
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
		),
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
			"Error Deleting AWS AppStream Application Fleet Association",
			fmt.Sprintf("Could not disassociate application %q from fleet %q: %v", applicationARN, fleetName, err),
		)
		return
	}
}
