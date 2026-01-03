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

	stackName := plan.StackName.ValueString()
	entitlementName := plan.EntitlementName.ValueString()
	applicationIdentifier := plan.ApplicationIdentifier.ValueString()

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.AssociateApplicationToEntitlement(
				ctx,
				&awsappstream.AssociateApplicationToEntitlementInput{
					StackName:             aws.String(stackName),
					EntitlementName:       aws.String(entitlementName),
					ApplicationIdentifier: aws.String(applicationIdentifier),
				},
			)
			return err
		},
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_AssociateApplicationToEntitlement.html
		util.WithRetryOnFns(
			util.IsEntitlementNotFoundException,
			util.IsOperationNotPermittedException,
			util.IsResourceNotFoundException,
		),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Application Entitlement Association",
			fmt.Sprintf("Could not associate application %q to entitlement %q (stack %q): %v",
				applicationIdentifier, entitlementName, stackName, err,
			),
		)
		return
	}

	newState, diags := r.readAssociateApplicationEntitlement(ctx, stackName, entitlementName, applicationIdentifier)
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
