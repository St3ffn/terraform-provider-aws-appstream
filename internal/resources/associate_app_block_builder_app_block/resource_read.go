// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_app_block_builder_app_block

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

	newState, diags := r.readAssociateAppBlockBuilderAppBlock(
		ctx, state.AppBlockBuilderName.ValueString(), state.AppBlockARN.ValueString(),
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

func (r *resource) readAssociateAppBlockBuilderAppBlock(
	ctx context.Context, appBlockBuilderName, appBlockARN string,
) (*model, diag.Diagnostics) {

	var diags diag.Diagnostics

	out, err := r.appstreamClient.DescribeAppBlockBuilderAppBlockAssociations(ctx, &awsappstream.DescribeAppBlockBuilderAppBlockAssociationsInput{
		AppBlockArn:         aws.String(appBlockARN),
		AppBlockBuilderName: aws.String(appBlockBuilderName),
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		diags.AddError(
			"Error Reading AWS AppStream App Block Builder App Block Association",
			fmt.Sprintf(
				"Could not read association of app block %q with app block builder %q: %v",
				appBlockARN, appBlockBuilderName, err,
			),
		)
		return nil, diags
	}

	if len(out.AppBlockBuilderAppBlockAssociations) == 0 {
		return nil, diags
	}

	state := &model{
		ID:                  types.StringValue(buildID(appBlockBuilderName, appBlockARN)),
		AppBlockBuilderName: types.StringValue(appBlockBuilderName),
		AppBlockARN:         types.StringValue(appBlockARN),
	}
	return state, diags

}
