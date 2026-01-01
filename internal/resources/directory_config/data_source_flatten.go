// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func flattenServiceAccountCredentials(
	ctx context.Context, creds *awstypes.ServiceAccountCredentials, diags *diag.Diagnostics,
) types.Object {

	if creds == nil {
		return types.ObjectNull(serviceAccountCredentialsObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		serviceAccountCredentialsObjectType.AttrTypes,
		serviceAccountCredentialsModel{
			AccountName:     util.StringOrNull(creds.AccountName),
			AccountPassword: util.StringOrNull(creds.AccountPassword),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(serviceAccountCredentialsObjectType.AttrTypes)
	}

	return obj
}

func flattenCertificateBasedAuthPropertiesData(
	ctx context.Context, props *awstypes.CertificateBasedAuthProperties, diags *diag.Diagnostics,
) types.Object {

	if props == nil {
		return types.ObjectNull(certificateBasedAuthPropertiesObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		certificateBasedAuthPropertiesObjectType.AttrTypes,
		certificateBasedAuthPropertiesModel{
			Status:                  types.StringValue(string(props.Status)),
			CertificateAuthorityARN: util.StringOrNull(props.CertificateAuthorityArn),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(certificateBasedAuthPropertiesObjectType.AttrTypes)
	}

	return obj
}
