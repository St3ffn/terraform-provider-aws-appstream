// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
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

	newState, diags := r.readAssociateUserStack(ctx, state)

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

func (r *resource) readAssociateUserStack(ctx context.Context, prior model) (*model, diag.Diagnostics) {
	var diags diag.Diagnostics

	stackName := prior.StackName.ValueString()
	userName := prior.UserName.ValueString()
	authenticationType := prior.AuthenticationType.ValueString()

	out, err := r.appstreamClient.DescribeUserStackAssociations(ctx, &awsappstream.DescribeUserStackAssociationsInput{
		AuthenticationType: awstypes.AuthenticationType(authenticationType),
		StackName:          aws.String(stackName),
		UserName:           aws.String(userName),
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		diags.AddError(
			"Error Reading AWS AppStream User Stack Association",
			fmt.Sprintf("Could not read association of user %q with stack %q: %v", userName, stackName, err),
		)
		return nil, diags
	}

	if len(out.UserStackAssociations) == 0 {
		return nil, diags
	}

	userStackAssociations := out.UserStackAssociations[0]

	state := &model{
		ID:                    types.StringValue(buildID(stackName, authenticationType, userName)),
		StackName:             types.StringValue(stackName),
		UserName:              types.StringValue(userName),
		AuthenticationType:    types.StringValue(authenticationType),
		SendEmailNotification: util.BoolOrNull(userStackAssociations.SendEmailNotification),
	}
	return state, diags
}
