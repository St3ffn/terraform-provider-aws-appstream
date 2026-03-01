// Copyright St3ffn 2025, 2026
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

func TestFlattenServiceAccountCredentials(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name  string
		prior types.Object
		aws   *awstypes.ServiceAccountCredentials
		want  types.Object
	}{
		{
			name:  "prior_null_returns_null",
			prior: types.ObjectNull(serviceAccountCredentialsObjectType.AttrTypes),
			aws:   &awstypes.ServiceAccountCredentials{AccountName: aws.String("svc")},
			want:  types.ObjectNull(serviceAccountCredentialsObjectType.AttrTypes),
		},
		{
			name:  "prior_unknown_returns_unknown",
			prior: types.ObjectUnknown(serviceAccountCredentialsObjectType.AttrTypes),
			aws:   &awstypes.ServiceAccountCredentials{AccountName: aws.String("svc")},
			want:  types.ObjectUnknown(serviceAccountCredentialsObjectType.AttrTypes),
		},
		{
			name: "aws_nil_preserves_account_name_and_clears_password",
			prior: mustObject(
				t,
				serviceAccountCredentialsObjectType.AttrTypes,
				resourceModelServiceAccountCredentials{
					AccountName:     types.StringValue("svc"),
					AccountPassword: types.StringValue("secret"),
				},
			),
			aws: nil,
			want: types.ObjectValueMust(
				serviceAccountCredentialsObjectType.AttrTypes,
				map[string]attr.Value{
					"account_name":     types.StringValue("svc"),
					"account_password": types.StringNull(),
				},
			),
		},
		{
			name: "normal_reconcile_password_is_cleared",
			prior: mustObject(
				t,
				serviceAccountCredentialsObjectType.AttrTypes,
				resourceModelServiceAccountCredentials{
					AccountName:     types.StringValue("svc"),
					AccountPassword: types.StringValue("secret"),
				},
			),
			aws: &awstypes.ServiceAccountCredentials{
				AccountName: aws.String("svc"),
			},
			want: types.ObjectValueMust(
				serviceAccountCredentialsObjectType.AttrTypes,
				map[string]attr.Value{
					"account_name":     types.StringValue("svc"),
					"account_password": types.StringNull(),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenServiceAccountCredentials(ctx, tt.prior, tt.aws, &diags)
			require.False(t, diags.HasError())

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenCertificateBasedAuthProperties(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.CertificateBasedAuthProperties
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(certificateBasedAuthPropertiesObjectType.AttrTypes),
		},
		{
			name: "status_only",
			in: &awstypes.CertificateBasedAuthProperties{
				Status: awstypes.CertificateBasedAuthStatusEnabled,
			},
			want: types.ObjectValueMust(
				certificateBasedAuthPropertiesObjectType.AttrTypes,
				map[string]attr.Value{
					"status":                    types.StringValue("ENABLED"),
					"certificate_authority_arn": types.StringNull(),
				},
			),
		},
		{
			name: "status_and_ca_arn",
			in: &awstypes.CertificateBasedAuthProperties{
				Status: awstypes.CertificateBasedAuthStatusEnabledNoDirectoryLoginFallback,
				CertificateAuthorityArn: aws.String(
					"arn:aws:acm-pca:eu-central-1:123456789012:certificate-authority/abcd",
				),
			},
			want: types.ObjectValueMust(
				certificateBasedAuthPropertiesObjectType.AttrTypes,
				map[string]attr.Value{
					"status": types.StringValue(
						"ENABLED_NO_DIRECTORY_LOGIN_FALLBACK",
					),
					"certificate_authority_arn": types.StringValue(
						"arn:aws:acm-pca:eu-central-1:123456789012:certificate-authority/abcd",
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenCertificateBasedAuthProperties(ctx, tt.in, &diags)

			require.False(t, diags.HasError())

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func mustObject[T any](t *testing.T, attrs map[string]attr.Type, in T) types.Object {
	t.Helper()

	obj, diags := types.ObjectValueFrom(context.Background(), attrs, in)
	require.False(t, diags.HasError(), "failed to build object value: %v", diags)
	return obj
}
