// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var entitlementAttributeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	},
}

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

	out, err := r.appstreamClient.DescribeEntitlements(ctx, &awsappstream.DescribeEntitlementsInput{
		StackName: aws.String(stackName),
		Name:      aws.String(name),
	})
	if err != nil {
		// respect cancellation/deadlines
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}
		if isAppStreamNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream Entitlement",
			fmt.Sprintf("Could not read entitlement %q in stack %q: %v", name, stackName, err),
		)
		return
	}

	if len(out.Entitlements) == 0 {
		// remove resource if missing
		resp.State.RemoveResource(ctx)
		return
	}

	e := out.Entitlements[0]
	if e.StackName == nil || e.Name == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var newState entitlementModel
	newState.ID = types.StringValue(buildEntitlementID(aws.ToString(e.StackName), aws.ToString(e.Name)))
	newState.StackName = types.StringValue(aws.ToString(e.StackName))
	newState.Name = types.StringValue(aws.ToString(e.Name))
	newState.AppVisibility = types.StringValue(string(e.AppVisibility))

	// optional description
	if e.Description != nil {
		newState.Description = types.StringValue(aws.ToString(e.Description))
	} else {
		newState.Description = types.StringNull()
	}

	newState.CreatedTime = stringFromTime(e.CreatedTime)
	newState.LastModifiedTime = stringFromTime(e.LastModifiedTime)

	newState.Attributes = flattenEntitlementAttributes(ctx, e.Attributes, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func flattenEntitlementAttributes(
	ctx context.Context, awsAttrs []awstypes.EntitlementAttribute, diags *diag.Diagnostics,
) types.Set {

	attrs := make([]entitlementAttributeModel, 0, len(awsAttrs))
	for _, a := range awsAttrs {
		attrs = append(attrs, entitlementAttributeModel{
			Name:  types.StringValue(aws.ToString(a.Name)),
			Value: types.StringValue(aws.ToString(a.Value)),
		})
	}

	setVal, d := types.SetValueFrom(ctx, entitlementAttributeObjectType, attrs)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(entitlementAttributeObjectType)
	}

	return setVal
}
