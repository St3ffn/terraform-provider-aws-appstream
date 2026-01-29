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

func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan model

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

	appBlockBuilderName := plan.AppBlockBuilderName.ValueString()
	appBlockARN := plan.AppBlockARN.ValueString()

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.AssociateAppBlockBuilderAppBlock(ctx, &awsappstream.AssociateAppBlockBuilderAppBlockInput{
				AppBlockArn:         aws.String(appBlockARN),
				AppBlockBuilderName: aws.String(appBlockBuilderName),
			})
			return err
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_AssociateAppBlockBuilderAppBlock.html
		util.WithRetryOnFns(
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
			util.IsResourceNotFoundException,
		),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream App Block Builder App Block Association",
			fmt.Sprintf(
				"Could not associate app block %q with app block builder %q: %v",
				appBlockARN, appBlockBuilderName, err,
			),
		)
		return
	}

	newState, diags := r.readAssociateAppBlockBuilderAppBlock(ctx, appBlockBuilderName, appBlockARN)
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
