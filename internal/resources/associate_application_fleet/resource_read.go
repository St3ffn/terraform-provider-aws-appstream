// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_application_fleet

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

func (r *resource) Read(ctx context.Context, req tfresource.ReadRequest, resp *tfresource.ReadResponse) {

	var state model

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

	newState, diags := r.readAssociateApplicationFleet(
		ctx, state.FleetName.ValueString(), state.ApplicationARN.ValueString(),
	)

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

func (r *resource) readAssociateApplicationFleet(
	ctx context.Context, fleetName, applicationARN string,
) (*model, diag.Diagnostics) {

	var diags diag.Diagnostics

	out, err := r.appstreamClient.DescribeApplicationFleetAssociations(ctx, &awsappstream.DescribeApplicationFleetAssociationsInput{
		FleetName:      aws.String(fleetName),
		ApplicationArn: aws.String(applicationARN),
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		diags.AddError(
			"Error Reading AWS AppStream Application Fleet Association",
			fmt.Sprintf(
				"Could not read association of application %q with fleet %q: %v",
				applicationARN, fleetName, err,
			),
		)
		return nil, diags
	}

	if len(out.ApplicationFleetAssociations) == 0 {
		return nil, diags
	}

	state := &model{
		ID:             types.StringValue(buildID(fleetName, applicationARN)),
		FleetName:      types.StringValue(fleetName),
		ApplicationARN: types.StringValue(applicationARN),
	}
	return state, diags

}
