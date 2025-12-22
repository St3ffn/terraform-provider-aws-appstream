// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *associateFleetStackResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {

	var state associateFleetStackModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	addAssociateFleetStackDiagnostics(state, &resp.Diagnostics, associateDiagnosticRead)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readAssociateFleetStack(ctx, state.FleetName.ValueString(), state.StackName.ValueString())

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if isContextCanceled(ctx) {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *associateFleetStackResource) readAssociateFleetStack(
	ctx context.Context, fleetName, stackName string,
) (*associateFleetStackModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	var nextToken *string

	for {
		out, err := r.appstreamClient.ListAssociatedStacks(ctx, &awsappstream.ListAssociatedStacksInput{
			FleetName: aws.String(fleetName),
			NextToken: nextToken,
		})
		if err != nil {
			if isContextCanceled(ctx) {
				return nil, diags
			}

			diags.AddError(
				"Error Reading AWS AppStream Fleet Stack Association",
				fmt.Sprintf(
					"Could not read association of fleet %q with stack %q: %v",
					fleetName, stackName, err,
				),
			)
			return nil, diags
		}

		for _, stack := range out.Names {
			if stack == stackName {
				state := &associateFleetStackModel{
					ID:        types.StringValue(buildAssociateFleetStackID(fleetName, stackName)),
					FleetName: types.StringValue(fleetName),
					StackName: types.StringValue(stackName),
				}
				return state, diags
			}
		}

		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		nextToken = out.NextToken
	}

	return nil, diags
}
