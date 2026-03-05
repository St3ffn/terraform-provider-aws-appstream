// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_app_block_builder_app_block

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

	addDiagnostics(state, &resp.Diagnostics, diagnosticDelete)
	if resp.Diagnostics.HasError() {
		return
	}

	appBlockBuilderName := state.AppBlockBuilderName.ValueString()
	appBlockARN := state.AppBlockARN.ValueString()

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.DisassociateAppBlockBuilderAppBlock(
				ctx,
				&awsappstream.DisassociateAppBlockBuilderAppBlockInput{
					AppBlockArn:         aws.String(appBlockARN),
					AppBlockBuilderName: aws.String(appBlockBuilderName),
				},
			)
			return err
		},
		util.WithTimeout(deleteRetryTimeout),
		util.WithInitBackoff(deleteRetryInitBackoff),
		util.WithMaxBackoff(deleteRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DisassociateAppBlockBuilderAppBlock.html
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
			"Error Deleting AWS AppStream App Block Builder App Block Association",
			fmt.Sprintf("Could not disassociate app block %q from app block builder %q: %v", appBlockARN, appBlockBuilderName, err),
		)
		return
	}
}
