// Copyright St3ffn 2025, 2026
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

func (r *resource) Update(ctx context.Context, req tfresource.UpdateRequest, resp *tfresource.UpdateResponse) {
	var plan resourceModel
	var state resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	diff := newResourceDiff(state, plan)
	imageBuilderARN := plan.ImageBuilderARN.ValueString()

	imageBuilderName, err := imageBuilderNameFromARN(imageBuilderARN)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			fmt.Sprintf("Could not parse image builder name from ARN %q: %v", imageBuilderARN, err),
		)
		return
	}

	added := make([]string, 0)
	removed := make([]string, 0)
	if diff.SoftwareNames.IsChanged() {
		added, removed = util.DiffStringSets(ctx, state.SoftwareNames, plan.SoftwareNames, &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	if len(added) > 0 || len(removed) > 0 || (diff.Deploy.IsChanged() && plan.Deploy.ValueBool()) {
		err = r.waitForImageBuilderAssociable(ctx, imageBuilderName)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Image Builder Software Association",
				fmt.Sprintf("Image builder %q did not reach RUNNING or STOPPED state in time: %v", imageBuilderARN, err),
			)
			return
		}
	}

	if len(removed) > 0 {
		err = util.RetryOn(
			ctx,
			func(ctx context.Context) error {
				_, err = r.appstreamClient.DisassociateSoftwareFromImageBuilder(ctx, &awsappstream.DisassociateSoftwareFromImageBuilderInput{
					ImageBuilderName: aws.String(imageBuilderName),
					SoftwareNames:    removed,
				})
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
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Image Builder Software Association",
				fmt.Sprintf("Could not disassociate software from image builder %q: %v", imageBuilderARN, err),
			)
			return
		}
	}

	if len(added) > 0 {
		err = util.RetryOn(
			ctx,
			func(ctx context.Context) error {
				_, err = r.appstreamClient.AssociateSoftwareToImageBuilder(ctx, &awsappstream.AssociateSoftwareToImageBuilderInput{
					ImageBuilderName: aws.String(imageBuilderName),
					SoftwareNames:    added,
				})
				return err
			},
			util.WithTimeout(retryTimeout),
			util.WithInitBackoff(retryInitBackoff),
			util.WithMaxBackoff(retryMaxBackoff),
			// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_AssociateSoftwareToImageBuilder.html
			util.WithRetryOnFns(
				util.IsConcurrentModificationException,
				util.IsOperationNotPermittedException,
			),
		)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating AWS AppStream Image Builder Software Association",
				fmt.Sprintf("Could not associate software to image builder %q: %v", imageBuilderARN, err),
			)
			return
		}
	}

	if diff.Deploy.IsChanged() && plan.Deploy.ValueBool() {
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
				"Error Updating AWS AppStream Image Builder Software Association",
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
