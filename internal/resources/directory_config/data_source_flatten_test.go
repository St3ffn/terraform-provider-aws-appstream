// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestFlattenServiceAccountCredentialsData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input *awstypes.ServiceAccountCredentials
		want  types.Object
	}{
		{
			name:  "nil_input_returns_null_object",
			input: nil,
			want:  types.ObjectNull(serviceAccountCredentialsObjectType.AttrTypes),
		},
		{
			name: "all_fields_set",
			input: &awstypes.ServiceAccountCredentials{
				AccountName:     aws.String("svc-user"),
				AccountPassword: aws.String("secret"),
			},
			want: types.ObjectValueMust(
				serviceAccountCredentialsObjectType.AttrTypes,
				map[string]attr.Value{
					"account_name":     types.StringValue("svc-user"),
					"account_password": types.StringValue("secret"),
				},
			),
		},
		{
			name: "optional_fields_nil",
			input: &awstypes.ServiceAccountCredentials{
				AccountName: aws.String("svc-user"),
			},
			want: types.ObjectValueMust(
				serviceAccountCredentialsObjectType.AttrTypes,
				map[string]attr.Value{
					"account_name":     types.StringValue("svc-user"),
					"account_password": types.StringNull(),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenServiceAccountCredentials(ctx, tt.input, &diags)
			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
			require.True(t, got.Equal(tt.want))
		})
	}
}

func TestFlattenCertificateBasedAuthPropertiesData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input *awstypes.CertificateBasedAuthProperties
		want  types.Object
	}{
		{
			name:  "nil_input_returns_null_object",
			input: nil,
			want:  types.ObjectNull(certificateBasedAuthPropertiesObjectType.AttrTypes),
		},
		{
			name: "all_fields_set",
			input: &awstypes.CertificateBasedAuthProperties{
				Status:                  awstypes.CertificateBasedAuthStatusEnabled,
				CertificateAuthorityArn: aws.String("arn:aws:acm-pca:us-east-1:123:certificate-authority/abc"),
			},
			want: types.ObjectValueMust(
				certificateBasedAuthPropertiesObjectType.AttrTypes,
				map[string]attr.Value{
					"status":                    types.StringValue("ENABLED"),
					"certificate_authority_arn": types.StringValue("arn:aws:acm-pca:us-east-1:123:certificate-authority/abc"),
				},
			),
		},
		{
			name: "arn_nil",
			input: &awstypes.CertificateBasedAuthProperties{
				Status: awstypes.CertificateBasedAuthStatusDisabled,
			},
			want: types.ObjectValueMust(
				certificateBasedAuthPropertiesObjectType.AttrTypes,
				map[string]attr.Value{
					"status":                    types.StringValue("DISABLED"),
					"certificate_authority_arn": types.StringNull(),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenCertificateBasedAuthPropertiesData(ctx, tt.input, &diags)
			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
			require.True(t, got.Equal(tt.want))
		})
	}
}
