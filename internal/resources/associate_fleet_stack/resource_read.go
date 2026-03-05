// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_fleet_stack

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Read fetches the current remote object and refreshes Terraform state from AWS.
// When the object no longer exists remotely, the resource is removed from state
// to converge Terraform with external deletions.
func (r *resource) Read(ctx context.Context, req tfresource.ReadRequest, resp *tfresource.ReadResponse) {

	var state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	addDiagnostics(state, &resp.Diagnostics, diagnosticRead)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readAssociateFleetStack(ctx, state.FleetName.ValueString(), state.StackName.ValueString())

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

func (r *resource) readAssociateFleetStack(
	ctx context.Context, fleetName, stackName string,
) (*resourceModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	var nextToken *string

	for {
		out, err := r.appstreamClient.ListAssociatedStacks(ctx, &awsappstream.ListAssociatedStacksInput{
			FleetName: aws.String(fleetName),
			NextToken: nextToken,
		})
		if err != nil {
			if util.IsContextCanceled(err) {
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
				state := &resourceModel{
					ID:        types.StringValue(buildID(fleetName, stackName)),
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
