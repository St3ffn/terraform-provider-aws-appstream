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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func TestFlattenServiceAccountCredentialsResource(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		prior         types.Object
		aws           *awstypes.ServiceAccountCredentials
		expectNull    bool
		expectUnknown bool
	}{
		{
			name:       "prior_null_returns_null",
			prior:      types.ObjectNull(serviceAccountCredentialsObjectType.AttrTypes),
			expectNull: true,
		},
		{
			name:          "prior_unknown_returns_unknown",
			prior:         types.ObjectUnknown(serviceAccountCredentialsObjectType.AttrTypes),
			expectUnknown: true,
		},
		{
			name: "aws_nil_returns_null",
			prior: mustObject(
				t,
				serviceAccountCredentialsObjectType.AttrTypes,
				serviceAccountCredentialsModel{
					AccountName:     types.StringValue("svc"),
					AccountPassword: types.StringValue("secret"),
				},
			),
			aws:        nil,
			expectNull: true,
		},
		{
			name: "normal_reconcile",
			prior: mustObject(
				t,
				serviceAccountCredentialsObjectType.AttrTypes,
				serviceAccountCredentialsModel{
					AccountName:     types.StringValue("svc"),
					AccountPassword: types.StringValue("secret"),
				},
			),
			aws: &awstypes.ServiceAccountCredentials{
				AccountName:     aws.String("svc"),
				AccountPassword: aws.String("secret"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			out := flattenServiceAccountCredentialsResource(ctx, tt.prior, tt.aws, &diags)
			require.False(t, diags.HasError())

			if tt.expectNull {
				require.True(t, out.IsNull())
				return
			}
			if tt.expectUnknown {
				require.True(t, out.IsUnknown())
				return
			}

			var model serviceAccountCredentialsModel
			diags = out.As(ctx, &model, basetypes.ObjectAsOptions{})
			require.False(t, diags.HasError())
			require.Equal(t, "svc", model.AccountName.ValueString())
		})
	}
}

func TestFlattenCertificateBasedAuthPropertiesResource(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		prior         types.Object
		aws           *awstypes.CertificateBasedAuthProperties
		expectNull    bool
		expectUnknown bool
	}{
		{
			name:       "prior_null_returns_null",
			prior:      types.ObjectNull(certificateBasedAuthPropertiesObjectType.AttrTypes),
			expectNull: true,
		},
		{
			name:          "prior_unknown_returns_unknown",
			prior:         types.ObjectUnknown(certificateBasedAuthPropertiesObjectType.AttrTypes),
			expectUnknown: true,
		},
		{
			name: "aws_nil_returns_null",
			prior: mustObject(
				t,
				certificateBasedAuthPropertiesObjectType.AttrTypes,
				certificateBasedAuthPropertiesModel{
					Status: types.StringValue("ENABLED"),
				},
			),
			aws:        nil,
			expectNull: true,
		},
		{
			name: "normal_reconcile",
			prior: mustObject(
				t,
				certificateBasedAuthPropertiesObjectType.AttrTypes,
				certificateBasedAuthPropertiesModel{
					Status: types.StringValue("DISABLED"),
				},
			),
			aws: &awstypes.CertificateBasedAuthProperties{
				Status: awstypes.CertificateBasedAuthStatusEnabled,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			out := flattenCertificateBasedAuthPropertiesResource(ctx, tt.prior, tt.aws, &diags)
			require.False(t, diags.HasError())

			if tt.expectNull {
				require.True(t, out.IsNull())
				return
			}
			if tt.expectUnknown {
				require.True(t, out.IsUnknown())
				return
			}

			var model certificateBasedAuthPropertiesModel
			diags = out.As(ctx, &model, basetypes.ObjectAsOptions{})
			require.False(t, diags.HasError())
			require.Equal(t, "ENABLED", model.Status.ValueString())
		})
	}
}

func mustObject[T any](t *testing.T, attrs map[string]attr.Type, in T) types.Object {
	t.Helper()

	obj, diags := types.ObjectValueFrom(context.Background(), attrs, in)
	require.False(t, diags.HasError(), "failed to build object value: %v", diags)
	return obj
}
