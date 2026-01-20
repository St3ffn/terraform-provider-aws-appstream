// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

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

	imageBuilderARN := state.ImageBuilderARN.ValueString()

	imageBuilderName, err := imageBuilderNameFromARN(imageBuilderARN)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			fmt.Sprintf("Could not parse image builder name from ARN %q: %v", imageBuilderARN, err),
		)
		return
	}

	input := &awsappstream.DisassociateSoftwareFromImageBuilderInput{
		ImageBuilderName: aws.String(imageBuilderName),
		SoftwareNames:    util.ExpandStringSetOrNil(ctx, state.SoftwareNames, &resp.Diagnostics),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	err = util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err = r.appstreamClient.DisassociateSoftwareFromImageBuilder(ctx, input)
			return err
		},
		util.WithTimeout(retryTimeout),
		util.WithInitBackoff(retryInitBackoff),
		util.WithMaxBackoff(retryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DisassociateSoftwareFromImageBuilder.html
		util.WithRetryOnFns(
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
		),
	)

	if err != nil {
		// if it's already gone, that's fine for delete.
		if util.IsAppStreamNotFound(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting AWS AppStream Image Builder Software Association",
			fmt.Sprintf("Could not disassociate software from image builder %q: %v", imageBuilderARN, err),
		)
		return
	}

	if state.Deploy.ValueBool() {
		err = util.RetryOn(
			ctx,
			func(ctx context.Context) error {
				_, err = r.appstreamClient.StartSoftwareDeploymentToImageBuilder(ctx, &awsappstream.StartSoftwareDeploymentToImageBuilderInput{
					ImageBuilderName:       aws.String(imageBuilderName),
					RetryFailedDeployments: aws.Bool(true),
				})
				return err
			},
			util.WithTimeout(retryTimeout),
			util.WithInitBackoff(retryInitBackoff),
			util.WithMaxBackoff(retryMaxBackoff),
			// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_StartSoftwareDeploymentToImageBuilder.html
			util.WithRetryOnFns(
				util.IsConcurrentModificationException,
				util.IsOperationNotPermittedException,
			),
		)

		if err != nil {
			// if it's already gone, that's fine for delete.
			if util.IsAppStreamNotFound(err) {
				return
			}

			resp.Diagnostics.AddError(
				"Error Deleting AWS AppStream Image Builder Software Association",
				fmt.Sprintf("Could not start software deployment to image builder %q: %v", imageBuilderARN, err),
			)
			return
		}
	}
}
