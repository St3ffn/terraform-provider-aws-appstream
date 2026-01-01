// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func expandServiceAccountCredentials(
	ctx context.Context, obj types.Object, diags *diag.Diagnostics,
) *awstypes.ServiceAccountCredentials {

	var m serviceAccountCredentialsModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	return &awstypes.ServiceAccountCredentials{
		AccountName:     util.StringPointerOrNil(m.AccountName),
		AccountPassword: util.StringPointerOrNil(m.AccountPassword),
	}
}

func expandCertificateBasedAuthProperties(
	ctx context.Context, obj types.Object, diags *diag.Diagnostics,
) *awstypes.CertificateBasedAuthProperties {

	var m certificateBasedAuthPropertiesModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	props := &awstypes.CertificateBasedAuthProperties{
		CertificateAuthorityArn: util.StringPointerOrNil(m.CertificateAuthorityARN),
	}

	if !m.Status.IsNull() && !m.Status.IsUnknown() {
		props.Status = awstypes.CertificateBasedAuthStatus(m.Status.ValueString())
	}

	if props.CertificateAuthorityArn == nil && props.Status == "" {
		return nil
	}

	return props
}
