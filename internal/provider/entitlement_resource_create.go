// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

	err := retryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.CreateEntitlement(ctx, &awsappstream.CreateEntitlementInput{
				StackName:     aws.String(stackName),
				Name:          aws.String(name),
				Description:   stringPointerOrNil(plan.Description),
				AppVisibility: awstypes.AppVisibility(plan.AppVisibility.ValueString()),
				Attributes:    awsAttrs,
			})
			return err
		},
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateEntitlement.html
		withRetryOnFns(
			isOperationNotPermittedException,
			isResourceNotFoundException),
	)

	if err != nil {
		if isContextCanceled(err) {
			return
		}

		if isResourceAlreadyExists(err) {
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

	newState, diags := r.readEntitlement(ctx, stackName, name)
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
