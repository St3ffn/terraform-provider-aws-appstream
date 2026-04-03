// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package stack

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

func TestFlattenStorageConnectorsResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name           string
		prior          types.Set
		aws            []awstypes.StorageConnector
		expectNull     bool
		expectUnknown  bool
		expectElements int
		expectRIDNull  bool
	}{
		{
			name:       "prior_null_returns_null",
			prior:      types.SetNull(storageConnectorObjectType),
			expectNull: true,
		},
		{
			name:          "prior_unknown_returns_unknown",
			prior:         types.SetUnknown(storageConnectorObjectType),
			expectUnknown: true,
		},
		{
			name: "aws_resource_identifier_not_adopted_when_not_configured",
			prior: mustSet(t, storageConnectorObjectType, []resourceModelStorageConnectors{
				{
					ConnectorType:              types.StringValue("HOMEFOLDERS"),
					ResourceIdentifier:         types.StringNull(),
					Domains:                    types.SetNull(types.StringType),
					DomainsRequireAdminConsent: types.SetNull(types.StringType),
				},
			}),
			aws: []awstypes.StorageConnector{
				{
					ConnectorType:      awstypes.StorageConnectorTypeHomefolders,
					ResourceIdentifier: aws.String("fs-123"),
				},
			},
			expectElements: 1,
			expectRIDNull:  true,
		},
		{
			name: "aws_connector_missing_returns_drifted_element",
			prior: mustSet(t, storageConnectorObjectType, []resourceModelStorageConnectors{
				{
					ConnectorType:              types.StringValue("HOMEFOLDERS"),
					ResourceIdentifier:         types.StringNull(),
					Domains:                    types.SetNull(types.StringType),
					DomainsRequireAdminConsent: types.SetNull(types.StringType),
				},
			}),
			expectElements: 1,
			expectRIDNull:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			out := flattenStorageConnectorsResource(ctx, tt.prior, tt.aws, &diags)
			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			if tt.expectNull {
				require.True(t, out.IsNull())
				return
			}
			if tt.expectUnknown {
				require.True(t, out.IsUnknown())
				return
			}

			var models []resourceModelStorageConnectors
			diags = out.ElementsAs(ctx, &models, false)
			require.False(t, diags.HasError())
			require.Len(t, models, tt.expectElements)

			if tt.expectRIDNull {
				require.True(t, models[0].ResourceIdentifier.IsNull())
			}
		})
	}
}

func TestFlattenUserSettingsResource_basic_mapping(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	prior := mustSet(t, userSettingObjectType, []resourceModelUserSettings{
		{
			Action:        types.StringValue("CLIPBOARD_COPY_FROM_LOCAL_DEVICE"),
			Permission:    types.StringValue("ENABLED"),
			MaximumLength: types.Int32Null(),
		},
	})

	awsSettings := []awstypes.UserSetting{
		{
			Action:     awstypes.ActionClipboardCopyFromLocalDevice,
			Permission: awstypes.PermissionEnabled,
		},
	}

	var diags diag.Diagnostics
	out := flattenUserSettingsResource(ctx, prior, awsSettings, &diags)
	require.False(t, diags.HasError())

	var models []resourceModelUserSettings
	diags = out.ElementsAs(ctx, &models, false)
	require.False(t, diags.HasError())
	require.Len(t, models, 1)

	require.Equal(t, "CLIPBOARD_COPY_FROM_LOCAL_DEVICE", models[0].Action.ValueString())
	require.Equal(t, "ENABLED", models[0].Permission.ValueString())
	require.True(t, models[0].MaximumLength.IsNull())
}

func TestFlattenApplicationSettingsResource_computed_bucket_preserved(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	prior := mustObject(t, applicationSettingsObjectType.AttrTypes, resourceModelApplicationSettings{
		Enabled:       types.BoolValue(true),
		SettingsGroup: types.StringNull(),
		S3BucketName:  types.StringNull(),
	})

	awsResp := &awstypes.ApplicationSettingsResponse{
		Enabled:       aws.Bool(false),
		SettingsGroup: aws.String("group"),
		S3BucketName:  aws.String("bucket"),
	}

	var diags diag.Diagnostics
	out := flattenApplicationSettingsResource(ctx, prior, awsResp, &diags)
	require.False(t, diags.HasError())

	var model resourceModelApplicationSettings
	diags = out.As(ctx, &model, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError())

	require.False(t, model.Enabled.ValueBool())
	require.True(t, model.SettingsGroup.IsNull())
	require.Equal(t, "bucket", model.S3BucketName.ValueString())
}

func TestFlattenContentRedirectionResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name          string
		prior         types.Object
		aws           *awstypes.ContentRedirection
		expectNull    bool
		expectUnknown bool
		want          *resourceModelContentRedirectionHostToClient
		wantNullHost  bool
	}{
		{
			name:       "prior_null_returns_null",
			prior:      types.ObjectNull(contentRedirectionObjectType.AttrTypes),
			aws:        &awstypes.ContentRedirection{},
			expectNull: true,
		},
		{
			name:          "prior_unknown_returns_unknown",
			prior:         types.ObjectUnknown(contentRedirectionObjectType.AttrTypes),
			aws:           &awstypes.ContentRedirection{},
			expectUnknown: true,
		},
		{
			name: "aws_nil_returns_null_host_to_client",
			prior: mustObject(t, contentRedirectionObjectType.AttrTypes, resourceModelContentRedirection{
				HostToClient: mustObject(t, contentRedirectionHostToClientObjectType.AttrTypes, resourceModelContentRedirectionHostToClient{
					Enabled:     types.BoolValue(true),
					AllowedUrls: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("https://example.com/*")}),
					DeniedUrls:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("https://example.com/admin/*")}),
				}),
			}),
			aws:          nil,
			wantNullHost: true,
		},
		{
			name: "basic_mapping",
			prior: mustObject(t, contentRedirectionObjectType.AttrTypes, resourceModelContentRedirection{
				HostToClient: mustObject(t, contentRedirectionHostToClientObjectType.AttrTypes, resourceModelContentRedirectionHostToClient{
					Enabled:     types.BoolValue(false),
					AllowedUrls: types.SetValueMust(types.StringType, []attr.Value{}),
					DeniedUrls:  types.SetValueMust(types.StringType, []attr.Value{}),
				}),
			}),
			aws: &awstypes.ContentRedirection{
				HostToClient: &awstypes.UrlRedirectionConfig{
					Enabled:     aws.Bool(true),
					AllowedUrls: []string{"https://example.com/*"},
					DeniedUrls:  []string{"https://example.com/admin/*"},
				},
			},
			want: &resourceModelContentRedirectionHostToClient{
				Enabled: types.BoolValue(true),
				AllowedUrls: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("https://example.com/*"),
				}),
				DeniedUrls: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("https://example.com/admin/*"),
				}),
			},
		},
		{
			name: "state_owned_urls_not_adopted_when_not_configured",
			prior: mustObject(t, contentRedirectionObjectType.AttrTypes, resourceModelContentRedirection{
				HostToClient: mustObject(t, contentRedirectionHostToClientObjectType.AttrTypes, resourceModelContentRedirectionHostToClient{
					Enabled:     types.BoolValue(false),
					AllowedUrls: types.SetNull(types.StringType),
					DeniedUrls:  types.SetNull(types.StringType),
				}),
			}),
			aws: &awstypes.ContentRedirection{
				HostToClient: &awstypes.UrlRedirectionConfig{
					Enabled:     aws.Bool(true),
					AllowedUrls: []string{"https://example.com/*"},
					DeniedUrls:  []string{"https://example.com/admin/*"},
				},
			},
			want: &resourceModelContentRedirectionHostToClient{
				Enabled:     types.BoolValue(true),
				AllowedUrls: types.SetNull(types.StringType),
				DeniedUrls:  types.SetNull(types.StringType),
			},
		},
		{
			name: "empty_aws_slices_become_empty_sets_when_owned",
			prior: mustObject(t, contentRedirectionObjectType.AttrTypes, resourceModelContentRedirection{
				HostToClient: mustObject(t, contentRedirectionHostToClientObjectType.AttrTypes, resourceModelContentRedirectionHostToClient{
					Enabled:     types.BoolValue(true),
					AllowedUrls: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("https://example.com/*")}),
					DeniedUrls:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("https://example.com/admin/*")}),
				}),
			}),
			aws: &awstypes.ContentRedirection{
				HostToClient: &awstypes.UrlRedirectionConfig{
					Enabled:     aws.Bool(false),
					AllowedUrls: []string{},
					DeniedUrls:  []string{},
				},
			},
			want: &resourceModelContentRedirectionHostToClient{
				Enabled:     types.BoolValue(false),
				AllowedUrls: types.SetValueMust(types.StringType, []attr.Value{}),
				DeniedUrls:  types.SetValueMust(types.StringType, []attr.Value{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			out := flattenContentRedirectionResource(ctx, tt.prior, tt.aws, &diags)
			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			if tt.expectNull {
				require.True(t, out.IsNull())
				return
			}

			if tt.expectUnknown {
				require.True(t, out.IsUnknown())
				return
			}

			var model resourceModelContentRedirection
			diags = out.As(ctx, &model, basetypes.ObjectAsOptions{})
			require.False(t, diags.HasError())

			if tt.wantNullHost {
				require.True(t, model.HostToClient.IsNull())
				return
			}

			var hostToClient resourceModelContentRedirectionHostToClient
			diags = model.HostToClient.As(ctx, &hostToClient, basetypes.ObjectAsOptions{})
			require.False(t, diags.HasError())
			require.Equal(t, *tt.want, hostToClient)
		})
	}
}

func TestFlattenAccessEndpointsResource_basic_mapping(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	prior := mustSet(t, accessEndpointObjectType, []resourceModelAccessEndpoints{
		{
			EndpointType: types.StringValue("STREAMING"),
			VpceID:       types.StringValue("vpce-123"),
		},
	})

	awsEndpoints := []awstypes.AccessEndpoint{
		{
			EndpointType: awstypes.AccessEndpointTypeStreaming,
			VpceId:       aws.String("vpce-123"),
		},
	}

	var diags diag.Diagnostics
	out := flattenAccessEndpointsResource(ctx, prior, awsEndpoints, &diags)
	require.False(t, diags.HasError())

	var models []resourceModelAccessEndpoints
	diags = out.ElementsAs(ctx, &models, false)
	require.False(t, diags.HasError())
	require.Len(t, models, 1)

	require.Equal(t, "STREAMING", models[0].EndpointType.ValueString())
	require.Equal(t, "vpce-123", models[0].VpceID.ValueString())
}

func TestFlattenStreamingExperienceSettingsResource_preferred_protocol_from_aws(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	prior := mustObject(
		t,
		streamingExperienceSettingsObjectType.AttrTypes,
		resourceModelStreamingExperienceSettings{
			PreferredProtocol: types.StringValue("TCP"),
		},
	)

	awsSettings := &awstypes.StreamingExperienceSettings{
		PreferredProtocol: awstypes.PreferredProtocolUdp,
	}

	var diags diag.Diagnostics
	out := flattenStreamingExperienceSettingsResource(ctx, prior, awsSettings, &diags)
	require.False(t, diags.HasError())

	var model resourceModelStreamingExperienceSettings
	diags = out.As(ctx, &model, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError())

	require.Equal(t, "UDP", model.PreferredProtocol.ValueString())
}

func mustSet[T any](t *testing.T, ot types.ObjectType, in []T) types.Set {
	t.Helper()

	set, diags := types.SetValueFrom(context.Background(), ot, in)
	require.False(t, diags.HasError(), "failed to build set value: %v", diags)
	return set
}

func mustObject[T any](t *testing.T, attrs map[string]attr.Type, in T) types.Object {
	t.Helper()

	obj, diags := types.ObjectValueFrom(context.Background(), attrs, in)
	require.False(t, diags.HasError(), "failed to build object value: %v", diags)
	return obj
}
