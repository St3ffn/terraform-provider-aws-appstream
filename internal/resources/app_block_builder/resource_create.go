// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

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
	var plan resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create app block builder because name must be known.",
		)
		return
	}

	name := plan.Name.ValueString()

	input := &awsappstream.CreateAppBlockBuilderInput{
		Name:         aws.String(name),
		InstanceType: aws.String(plan.InstanceType.ValueString()),
		Platform:     awstypes.AppBlockBuilderPlatformType(plan.Platform.ValueString()),
		VpcConfig:    expandVPCConfig(ctx, plan.VPCConfig, &resp.Diagnostics),
	}

	input.Description = util.StringPointerOrNil(plan.Description)
	input.DisplayName = util.StringPointerOrNil(plan.DisplayName)
	input.DisableIMDSV1 = util.BoolPointerOrNil(plan.DisableIMDSV1)

	input.IamRoleArn = util.StringPointerOrNil(plan.IAMRoleARN)
	input.EnableDefaultInternetAccess = util.BoolPointerOrNil(plan.EnableDefaultInternetAccess)

	if !plan.AccessEndpoints.IsNull() && !plan.AccessEndpoints.IsUnknown() {
		input.AccessEndpoints = expandAccessEndpoints(ctx, plan.AccessEndpoints, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	out, err := util.RetryOnValue(
		ctx,
		func(ctx context.Context) (*awsappstream.CreateAppBlockBuilderOutput, error) {
			return r.appstreamClient.CreateAppBlockBuilder(ctx, input)
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateAppBlockBuilder.html
		util.WithRetryOnFns(
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
			util.IsResourceNotAvailableException,
			util.IsResourceNotFoundException,
			util.IsInvalidRoleException,
		),
	)

	if err != nil {
		if util.IsResourceAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream App Block Builder Already Exists",
				fmt.Sprintf(
					"A app block builder named %q already exists. To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> %q",
					name, name,
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream App Block Builder",
			fmt.Sprintf("Could not create app block builder %q: %v", name, err),
		)
		return
	}

	if out.AppBlockBuilder != nil && out.AppBlockBuilder.Arn != nil {
		_, tagDiags := r.tags.Apply(ctx, aws.ToString(out.AppBlockBuilder.Arn), plan.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	newState, diags := r.readAppBlockBuilder(ctx, plan)
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
