// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan model

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

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.CreateEntitlement(ctx, &awsappstream.CreateEntitlementInput{
				StackName:     aws.String(stackName),
				Name:          aws.String(name),
				Description:   util.StringPointerOrNil(plan.Description),
				AppVisibility: awstypes.AppVisibility(plan.AppVisibility.ValueString()),
				Attributes:    awsAttrs,
			})
			return err
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateEntitlement.html
		util.WithRetryOnFns(
			util.IsOperationNotPermittedException,
			util.IsResourceNotFoundException),
	)

	if err != nil {
		if util.IsResourceAlreadyExists(err) || util.IsEntitlementAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Entitlement Already Exists",
				fmt.Sprintf(
					"An entitlement named %q already exists in stack %q. "+
						"To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> %q",
					name, stackName, buildID(stackName, name),
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

	newState, diags := r.readEntitlement(ctx, plan)
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
