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

func (r *entitlementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entitlementModel

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
				"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
		)
		return
	}

	stackName := state.StackName.ValueString()
	name := state.Name.ValueString()

	newState, diags := r.readEntitlement(ctx, stackName, name)
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

func (r *entitlementResource) readEntitlement(
	ctx context.Context, stackName, name string,
) (*entitlementModel, diag.Diagnostics) {

	var diags diag.Diagnostics

	out, err := r.appstreamClient.DescribeEntitlements(ctx, &awsappstream.DescribeEntitlementsInput{
		StackName: aws.String(stackName),
		Name:      aws.String(name),
	})
	if err != nil {
		if isContextCanceled(ctx) {
			return nil, diags
		}

		if isAppStreamNotFound(err) {
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

	state := &entitlementModel{
		ID:               types.StringValue(buildEntitlementID(stackName, name)),
		StackName:        types.StringValue(aws.ToString(e.StackName)),
		Name:             types.StringValue(aws.ToString(e.Name)),
		Description:      stringOrNull(e.Description),
		AppVisibility:    types.StringValue(string(e.AppVisibility)),
		CreatedTime:      stringFromTime(e.CreatedTime),
		LastModifiedTime: stringFromTime(e.LastModifiedTime),
		Attributes:       flattenEntitlementAttributes(ctx, e.Attributes, &diags),
	}

	if diags.HasError() {
		return nil, diags
	}

	return state, diags
}
