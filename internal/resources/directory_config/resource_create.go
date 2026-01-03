// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
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

	// Validate required attributes
	if plan.DirectoryName.IsNull() || plan.DirectoryName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot create directory config because directory_name must be known.",
		)
		return
	}

	name := plan.DirectoryName.ValueString()

	input := &awsappstream.CreateDirectoryConfigInput{
		DirectoryName: aws.String(name),
		OrganizationalUnitDistinguishedNames: util.ExpandStringSetOrNil(
			ctx, plan.OrganizationalUnitDistinguishedNames, &resp.Diagnostics,
		),
	}

	if !plan.ServiceAccountCredentials.IsNull() && !plan.ServiceAccountCredentials.IsUnknown() {
		input.ServiceAccountCredentials = expandServiceAccountCredentials(
			ctx, plan.ServiceAccountCredentials, &resp.Diagnostics,
		)
	}

	if !plan.CertificateBasedAuthProperties.IsNull() && !plan.CertificateBasedAuthProperties.IsUnknown() {
		input.CertificateBasedAuthProperties = expandCertificateBasedAuthProperties(
			ctx, plan.CertificateBasedAuthProperties, &resp.Diagnostics,
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	err := util.RetryOn(
		ctx,
		func(ctx context.Context) error {
			_, err := r.appstreamClient.CreateDirectoryConfig(ctx, input)
			return err
		},
		util.WithTimeout(createRetryTimeout),
		util.WithInitBackoff(createRetryInitBackoff),
		util.WithMaxBackoff(createRetryMaxBackoff),
		// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_CreateDirectoryConfig.html
		util.WithRetryOnFns(
			util.IsOperationNotPermittedException,
		),
	)

	if err != nil {
		if util.IsResourceAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Directory Config Already Exists",
				fmt.Sprintf(
					"An directory config named %q already exists. "+
						"To manage it with Terraform, import it using:\n\n"+
						"  terraform import <resource_address> <directory_name>",
					name,
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS Directory Config",
			fmt.Sprintf("Could not create directory config %q: %v", name, err),
		)
		return
	}

	newState, diags := r.readDirectoryConfig(ctx, plan)
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
