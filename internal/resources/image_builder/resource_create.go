// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
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
			"Cannot create image builder because name must be known.",
		)
		return
	}

	name := plan.Name.ValueString()

	input := &awsappstream.CreateImageBuilderInput{
		Name:         aws.String(name),
		InstanceType: aws.String(plan.InstanceType.ValueString()),
	}

	input.ImageName = util.StringPointerOrNil(plan.ImageName)
	input.ImageArn = util.StringPointerOrNil(plan.ImageARN)
	input.Description = util.StringPointerOrNil(plan.Description)
	input.DisplayName = util.StringPointerOrNil(plan.DisplayName)

	if !plan.VPCConfig.IsNull() && !plan.VPCConfig.IsUnknown() {
		input.VpcConfig = expandVPCConfig(ctx, plan.VPCConfig, &resp.Diagnostics)
	}

	input.IamRoleArn = util.StringPointerOrNil(plan.IAMRoleARN)
	input.EnableDefaultInternetAccess = util.BoolPointerOrNil(plan.EnableDefaultInternetAccess)

	if !plan.DomainJoinInfo.IsNull() && !plan.DomainJoinInfo.IsUnknown() {
		input.DomainJoinInfo = expandDomainJoinInfo(ctx, plan.DomainJoinInfo, &resp.Diagnostics)
	}

	input.AppstreamAgentVersion = util.StringPointerOrNil(plan.AppstreamAgentVersion)

	if !plan.AccessEndpoints.IsNull() && !plan.AccessEndpoints.IsUnknown() {
		input.AccessEndpoints = expandAccessEndpoints(ctx, plan.AccessEndpoints, &resp.Diagnostics)
	}

	if !plan.RootVolumeConfig.IsNull() && !plan.RootVolumeConfig.IsUnknown() {
		input.RootVolumeConfig = expandRootVolumeConfig(ctx, plan.RootVolumeConfig, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var out *awsappstream.CreateImageBuilderOutput
	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			var err error
			out, err = r.appstreamClient.CreateImageBuilder(ctx, input)
			return err
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateImageBuilder.html
		util.WithRetryOnFns(
			util.IsConcurrentModificationException,
			util.IsOperationNotPermittedException,
			util.IsResourceNotAvailableException,
			util.IsResourceNotFoundException,
		),
	)

	if err != nil {
		if util.IsResourceAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Image Builder Already Exists",
				fmt.Sprintf(
					"A image builder named %q already exists. To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> %q",
					name, name,
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Image Builder",
			fmt.Sprintf("Could not create image builder %q: %v", name, err),
		)
		return
	}

	if out.ImageBuilder != nil && out.ImageBuilder.Arn != nil {
		_, tagDiags := r.tags.Apply(ctx, aws.ToString(out.ImageBuilder.Arn), plan.Tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	newState, diags := r.readImageBuilder(ctx, plan)
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
