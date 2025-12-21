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

func (r *associateApplicationEntitlementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan associateApplicationEntitlementModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	addAssocPartsDiagnostics(plan, &resp.Diagnostics, assocDiagPlan)
	if resp.Diagnostics.HasError() {
		return
	}

	stackName := plan.StackName.ValueString()
	entitlementName := plan.EntitlementName.ValueString()
	applicationIdentifier := plan.ApplicationIdentifier.ValueString()

	_, err := r.appstreamClient.AssociateApplicationToEntitlement(ctx, &awsappstream.AssociateApplicationToEntitlementInput{
		StackName:             aws.String(stackName),
		EntitlementName:       aws.String(entitlementName),
		ApplicationIdentifier: aws.String(applicationIdentifier),
	})
	if err != nil {
		if isContextCanceled(ctx) {
			return
		}

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
