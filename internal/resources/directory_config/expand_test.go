// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package directory_config

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExpandServiceAccountCredentials(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		obj  types.Object
		want *awstypes.ServiceAccountCredentials
	}{
		{
			name: "both_fields_set",
			obj: types.ObjectValueMust(
				serviceAccountCredentialsObjectType.AttrTypes,
				map[string]attr.Value{
					"account_name":     types.StringValue("svc"),
					"account_password": types.StringValue("secret"),
				},
			),
			want: &awstypes.ServiceAccountCredentials{
				AccountName:     aws.String("svc"),
				AccountPassword: aws.String("secret"),
			},
		},
		{
			name: "password_null_still_expands",
			obj: types.ObjectValueMust(
				serviceAccountCredentialsObjectType.AttrTypes,
				map[string]attr.Value{
					"account_name":     types.StringValue("svc"),
					"account_password": types.StringNull(),
				},
			),
			want: &awstypes.ServiceAccountCredentials{
				AccountName: aws.String("svc"),
			},
		},
		{
			name: "both_null_returns_struct_with_nil_fields",
			obj: types.ObjectValueMust(
				serviceAccountCredentialsObjectType.AttrTypes,
				map[string]attr.Value{
					"account_name":     types.StringNull(),
					"account_password": types.StringNull(),
				},
			),
			want: &awstypes.ServiceAccountCredentials{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := expandServiceAccountCredentials(ctx, tt.obj, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExpandCertificateBasedAuthProperties(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		obj  types.Object
		want *awstypes.CertificateBasedAuthProperties
	}{
		{
			name: "status_only",
			obj: types.ObjectValueMust(
				certificateBasedAuthPropertiesObjectType.AttrTypes,
				map[string]attr.Value{
					"status":                    types.StringValue("ENABLED"),
					"certificate_authority_arn": types.StringNull(),
				},
			),
			want: &awstypes.CertificateBasedAuthProperties{
				Status: awstypes.CertificateBasedAuthStatusEnabled,
			},
		},
		{
			name: "arn_only",
			obj: types.ObjectValueMust(
				certificateBasedAuthPropertiesObjectType.AttrTypes,
				map[string]attr.Value{
					"status":                    types.StringNull(),
					"certificate_authority_arn": types.StringValue("arn:aws:acm-pca:eu-central-1:123:certificate-authority/abc"),
				},
			),
			want: &awstypes.CertificateBasedAuthProperties{
				CertificateAuthorityArn: aws.String("arn:aws:acm-pca:eu-central-1:123:certificate-authority/abc"),
			},
		},
		{
			name: "both_fields_set",
			obj: types.ObjectValueMust(
				certificateBasedAuthPropertiesObjectType.AttrTypes,
				map[string]attr.Value{
					"status":                    types.StringValue("DISABLED"),
					"certificate_authority_arn": types.StringValue("arn:aws:acm-pca:eu-central-1:123:certificate-authority/abc"),
				},
			),
			want: &awstypes.CertificateBasedAuthProperties{
				Status:                  awstypes.CertificateBasedAuthStatusDisabled,
				CertificateAuthorityArn: aws.String("arn:aws:acm-pca:eu-central-1:123:certificate-authority/abc"),
			},
		},
		{
			name: "both_null_returns_nil",
			obj: types.ObjectValueMust(
				certificateBasedAuthPropertiesObjectType.AttrTypes,
				map[string]attr.Value{
					"status":                    types.StringNull(),
					"certificate_authority_arn": types.StringNull(),
				},
			),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := expandCertificateBasedAuthProperties(ctx, tt.obj, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
