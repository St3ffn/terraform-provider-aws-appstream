// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (ds *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config model

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	if config.DirectoryName.IsNull() || config.DirectoryName.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Terraform Configuration",
			"Cannot read directory config because directory_name must be set and known.",
		)
		return
	}

	name := config.DirectoryName.ValueString()

	out, err := ds.appstreamClient.DescribeDirectoryConfigs(ctx, &awsappstream.DescribeDirectoryConfigsInput{
		DirectoryNames: []string{name},
	})
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Directory Config Not Found",
				fmt.Sprintf("No directory config named %q was found.", name),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream Directory Config",
			fmt.Sprintf("Could not read directory config %q: %v", name, err),
		)
		return
	}

	if len(out.DirectoryConfigs) == 0 {
		resp.Diagnostics.AddError(
			"AWS AppStream Directory Config Not Found",
			fmt.Sprintf("No directory config named %q was found.", name),
		)
		return
	}

	directoryConfig := out.DirectoryConfigs[0]
	if directoryConfig.DirectoryName == nil {
		resp.Diagnostics.AddError(
			"Unexpected AWS Response",
			fmt.Sprintf("Directory config %q was returned without required identifiers.", name),
		)
		return
	}

	state := &model{
		ID:            types.StringValue(aws.ToString(directoryConfig.DirectoryName)),
		DirectoryName: types.StringValue(aws.ToString(directoryConfig.DirectoryName)),
		OrganizationalUnitDistinguishedNames: util.SetStringOrNull(
			ctx,
			directoryConfig.OrganizationalUnitDistinguishedNames,
			&resp.Diagnostics,
		),
		ServiceAccountCredentials: flattenServiceAccountCredentials(
			ctx,
			directoryConfig.ServiceAccountCredentials,
			&resp.Diagnostics,
		),
		CertificateBasedAuthProperties: flattenCertificateBasedAuthPropertiesData(
			ctx,
			directoryConfig.CertificateBasedAuthProperties,
			&resp.Diagnostics,
		),
		CreatedTime: util.StringFromTime(directoryConfig.CreatedTime),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
