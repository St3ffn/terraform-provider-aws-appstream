// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Create validates plan input, performs the required AWS create/associate calls,
// retries configured transient API errors, and then reads the remote object back
// so state is written from the authoritative AWS response.
func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	addDiagnostics(plan, &resp.Diagnostics, diagnosticPlan)
	if resp.Diagnostics.HasError() {
		return
	}

	imageBuilderARN := plan.ImageBuilderARN.ValueString()

	imageBuilderName, err := imageBuilderNameFromARN(imageBuilderARN)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			fmt.Sprintf("Could not parse image builder name from ARN %q: %v", imageBuilderARN, err),
		)
		return
	}

	err = r.waitForImageBuilderAssociable(ctx, imageBuilderName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Image Builder Software Association",
			fmt.Sprintf("Image builder %q did not reach RUNNING or STOPPED state in time: %v", imageBuilderARN, err),
		)
		return
	}

	input := &awsappstream.AssociateSoftwareToImageBuilderInput{
		ImageBuilderName: aws.String(imageBuilderName),
		SoftwareNames:    util.ExpandStringSetOrNil(ctx, plan.SoftwareNames, &resp.Diagnostics),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	err = util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err = r.appstreamClient.AssociateSoftwareToImageBuilder(ctx, input)
			return err
		},
		util.WithTimeout(retryTimeout),
		util.WithInitBackoff(retryInitBackoff),
		util.WithMaxBackoff(retryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_AssociateSoftwareToImageBuilder.html
		util.WithRetryOnFns(
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
			util.IsResourceNotFoundException,
		),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Image Builder Software Association",
			fmt.Sprintf("Could not associate software to image builder %q: %v", imageBuilderARN, err),
		)
		return
	}

	if plan.Deploy.ValueBool() {
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
			resp.Diagnostics.AddError(
				"Error Creating AWS AppStream Image Builder Software Association",
				fmt.Sprintf("Could not start software deployment to image builder %q: %v", imageBuilderARN, err),
			)
			return
		}
	}

	newState, diags := r.readAssociationImageBuilderSoftware(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if ctx.Err() != nil {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func imageBuilderNameFromARN(arn string) (string, error) {
	// expected arn:aws:appstream:<region>:<account>:image-builder/<name>
	const prefix = "image-builder/"
	idx := strings.LastIndex(arn, prefix)
	if idx == -1 || idx+len(prefix) >= len(arn) {
		return "", fmt.Errorf("invalid image builder ARN format")
	}

	return arn[idx+len(prefix):], nil
}
