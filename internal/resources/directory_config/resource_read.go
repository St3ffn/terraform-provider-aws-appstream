// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

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

	if state.DirectoryName.IsNull() || state.DirectoryName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform State",
			"Required attribute directory_name is missing from state. "+
				"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
		)
		return
	}

	newState, diags := r.readDirectoryConfig(ctx, state)
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

func (r *resource) readDirectoryConfig(ctx context.Context, prior model) (*model, diag.Diagnostics) {
	var diags diag.Diagnostics

	name := prior.DirectoryName.ValueString()

	out, err := r.appstreamClient.DescribeDirectoryConfigs(ctx, &awsappstream.DescribeDirectoryConfigsInput{
		DirectoryNames: []string{name},
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		if util.IsAppStreamNotFound(err) {
			return nil, diags
		}

		diags.AddError(
			"Error Reading AWS AppStream Directory Config",
			fmt.Sprintf("Could not read directory config %q: %v", name, err),
		)
		return nil, diags
	}

	if len(out.DirectoryConfigs) == 0 {
		return nil, diags
	}

	directoryConfig := out.DirectoryConfigs[0]
	if directoryConfig.DirectoryName == nil {
		return nil, diags
	}

	state := &model{
		ID:                                   types.StringValue(aws.ToString(directoryConfig.DirectoryName)),
		DirectoryName:                        types.StringValue(aws.ToString(directoryConfig.DirectoryName)),
		OrganizationalUnitDistinguishedNames: util.SetStringOrNull(ctx, directoryConfig.OrganizationalUnitDistinguishedNames, &diags),
		ServiceAccountCredentials: flattenServiceAccountCredentialsResource(
			ctx, prior.ServiceAccountCredentials, directoryConfig.ServiceAccountCredentials, &diags,
		),
		CertificateBasedAuthProperties: flattenCertificateBasedAuthPropertiesResource(
			ctx, prior.CertificateBasedAuthProperties, directoryConfig.CertificateBasedAuthProperties, &diags,
		),
		CreatedTime: util.StringFromTime(directoryConfig.CreatedTime),
	}

	if diags.HasError() {
		return nil, diags
	}

	return state, diags
}
