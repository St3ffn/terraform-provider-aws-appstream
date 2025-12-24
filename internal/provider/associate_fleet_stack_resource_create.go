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

func (r *associateFleetStackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan associateFleetStackModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	addAssociateFleetStackDiagnostics(plan, &resp.Diagnostics, associateDiagnosticPlan)
	if resp.Diagnostics.HasError() {
		return
	}

	fleetName := plan.FleetName.ValueString()
	stackName := plan.StackName.ValueString()

	err := retryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.AssociateFleet(ctx, &awsappstream.AssociateFleetInput{
				FleetName: aws.String(fleetName),
				StackName: aws.String(stackName),
			})
			return err
		},
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_AssociateFleet.html
		withRetryOnFns(
			isOperationNotPermittedException,
			isResourceNotFoundException,
			isConcurrentModificationException,
		),
	)

	if err != nil {
		if isContextCanceled(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Fleet Stack Association",
			fmt.Sprintf("Could not associate fleet %q to stack %q: %v",
				fleetName, stackName, err,
			),
		)
		return
	}

	newState, diags := r.readAssociateFleetStack(ctx, fleetName, stackName)
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
