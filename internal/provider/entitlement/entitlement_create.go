// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *entitlementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entitlementModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	if plan.StackName.IsNull() || plan.StackName.IsUnknown() ||
		plan.Name.IsNull() || plan.Name.IsUnknown() ||
		plan.AppVisibility.IsNull() || plan.AppVisibility.IsUnknown() ||
		plan.Attributes.IsNull() || plan.Attributes.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create entitlement because stack_name, name, app_visibility, and attributes must be known.",
		)
		return
	}

	stackName := plan.StackName.ValueString()
	name := plan.Name.ValueString()

	awsAttrs := expandEntitlementAttributes(ctx, plan.Attributes, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var description *string
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		description = plan.Description.ValueStringPointer()
	}

	_, err := r.appStreamClient.CreateEntitlement(ctx, &awsappstream.CreateEntitlementInput{
		StackName:     aws.String(stackName),
		Name:          aws.String(name),
		Description:   description,
		AppVisibility: awstypes.AppVisibility(plan.AppVisibility.ValueString()),
		Attributes:    awsAttrs,
	})
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}

		if isAppStreamAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Entitlement Already Exists",
				fmt.Sprintf(
					"An entitlement named %q already exists in stack %q. "+
						"To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> %q",
					name, stackName, buildEntitlementID(stackName, name),
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Entitlement",
			fmt.Sprintf("Could not create entitlement %q in stack %q: %v", name, stackName, err),
		)
		return
	}

	plan.ID = types.StringValue(buildEntitlementID(stackName, name))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func buildEntitlementID(stackName, name string) string {
	return fmt.Sprintf("%s|%s", stackName, name)
}

func expandEntitlementAttributes(
	ctx context.Context, tfAttrs types.Set, diags *diag.Diagnostics,
) []awstypes.EntitlementAttribute {

	var attrs []entitlementAttributeModel
	diags.Append(tfAttrs.ElementsAs(ctx, &attrs, false)...)
	if diags.HasError() {
		return nil
	}

	awsAttrs := make([]awstypes.EntitlementAttribute, 0, len(attrs))
	for _, a := range attrs {
		if a.Name.IsNull() || a.Name.IsUnknown() ||
			a.Value.IsNull() || a.Value.IsUnknown() {

			diags.AddError(
				"Invalid Terraform Plan",
				"All entitlement attributes must have known, non-null name and value.",
			)
			return nil
		}

		awsAttrs = append(awsAttrs, awstypes.EntitlementAttribute{
			Name:  aws.String(a.Name.ValueString()),
			Value: aws.String(a.Value.ValueString()),
		})
	}

	return awsAttrs
}
