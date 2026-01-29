// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package directory_config

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var serviceAccountCredentialsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"account_name":     types.StringType,
		"account_password": types.StringType,
	},
}

var certificateBasedAuthPropertiesObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"status":                    types.StringType,
		"certificate_authority_arn": types.StringType,
	},
}

func flattenServiceAccountCredentials(
	ctx context.Context, prior types.Object, awsCreds *awstypes.ServiceAccountCredentials, diags *diag.Diagnostics,
) types.Object {

	// user never managed it
	if prior.IsNull() {
		return types.ObjectNull(serviceAccountCredentialsObjectType.AttrTypes)
	}

	// planning phase
	if prior.IsUnknown() {
		return types.ObjectUnknown(serviceAccountCredentialsObjectType.AttrTypes)
	}

	// drift or redaction: block disappears
	if awsCreds == nil {
		return types.ObjectNull(serviceAccountCredentialsObjectType.AttrTypes)
	}

	var priorProps serviceAccountCredentialsModel
	diags.Append(prior.As(ctx, &priorProps, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return types.ObjectNull(serviceAccountCredentialsObjectType.AttrTypes)
	}

	// normal reconcile
	obj, d := types.ObjectValueFrom(
		ctx,
		serviceAccountCredentialsObjectType.AttrTypes,
		serviceAccountCredentialsModel{
			AccountName:     util.StringOrNull(awsCreds.AccountName),
			AccountPassword: priorProps.AccountPassword,
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(serviceAccountCredentialsObjectType.AttrTypes)
	}

	return obj
}

func flattenCertificateBasedAuthProperties(
	ctx context.Context, awsProps *awstypes.CertificateBasedAuthProperties, diags *diag.Diagnostics,
) types.Object {

	if awsProps == nil {
		return types.ObjectNull(certificateBasedAuthPropertiesObjectType.AttrTypes)
	}

	obj, d := types.ObjectValueFrom(
		ctx,
		certificateBasedAuthPropertiesObjectType.AttrTypes,
		certificateBasedAuthPropertiesModel{
			Status:                  types.StringValue(string(awsProps.Status)),
			CertificateAuthorityARN: util.StringOrNull(awsProps.CertificateAuthorityArn),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(certificateBasedAuthPropertiesObjectType.AttrTypes)
	}

	return obj
}
