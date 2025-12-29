// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

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

	if state.StackName.IsNull() || state.StackName.IsUnknown() ||
		state.Name.IsNull() || state.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Required attributes stack_name and name are missing from state. "+
				"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the tfresource.",
		)
		return
	}

	newState, diags := r.readEntitlement(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if util.IsContextCanceled(ctx.Err()) {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *resource) readEntitlement(ctx context.Context, prior model) (*model, diag.Diagnostics) {

	var diags diag.Diagnostics

	stackName := prior.StackName.ValueString()
	name := prior.Name.ValueString()

	out, err := r.appstreamClient.DescribeEntitlements(ctx, &awsappstream.DescribeEntitlementsInput{
		StackName: aws.String(stackName),
		Name:      aws.String(name),
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		if util.IsAppStreamNotFound(err) {
			return nil, diags
		}

		diags.AddError(
			"Error Reading AWS AppStream Entitlement",
			fmt.Sprintf("Could not read entitlement %q in stack %q: %v", name, stackName, err),
		)
		return nil, diags
	}

	if len(out.Entitlements) == 0 {
		return nil, diags
	}

	e := out.Entitlements[0]
	if e.StackName == nil || e.Name == nil {
		return nil, diags
	}

	state := &model{
		ID:            types.StringValue(buildID(stackName, name)),
		StackName:     types.StringValue(aws.ToString(e.StackName)),
		Name:          types.StringValue(aws.ToString(e.Name)),
		Description:   util.StringOrNull(e.Description),
		AppVisibility: types.StringValue(string(e.AppVisibility)),
		CreatedTime:   util.StringFromTime(e.CreatedTime),
		Attributes:    flattenAttributes(ctx, e.Attributes, &diags),
	}

	if diags.HasError() {
		return nil, diags
	}

	return state, diags
}
