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

func (r *resource) Update(ctx context.Context, req tfresource.UpdateRequest, resp *tfresource.UpdateResponse) {
	var plan model
	var state model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if plan.DirectoryName.IsNull() || plan.DirectoryName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Plan",
			"Cannot update directory config because directory_name must be known.",
		)
		return
	}

	name := plan.DirectoryName.ValueString()

	// guard against unexpected identity drift
	if !state.DirectoryName.IsNull() && !state.DirectoryName.IsUnknown() {
		if state.DirectoryName.ValueString() != name {
			resp.Diagnostics.AddError(
				"Unexpected Update Request",
				"Directory config identity (directory_name) changed during update. This should trigger replacement. Please report this issue.",
			)
			return
		}
	}

	input := &awsappstream.UpdateDirectoryConfigInput{
		DirectoryName: aws.String(name),
		OrganizationalUnitDistinguishedNames: util.ExpandStringSetOrNil(
			ctx, plan.OrganizationalUnitDistinguishedNames, &resp.Diagnostics,
		),
		ServiceAccountCredentials: expandServiceAccountCredentials(
			ctx, plan.ServiceAccountCredentials, &resp.Diagnostics,
		),
		CertificateBasedAuthProperties: expandCertificateBasedAuthProperties(
			ctx, plan.CertificateBasedAuthProperties, &resp.Diagnostics,
		),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.appstreamClient.UpdateDirectoryConfig(ctx, input)
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Updating AWS AppStream Directory Config",
			fmt.Sprintf("Could not update directory config %q: %v", name, err),
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
