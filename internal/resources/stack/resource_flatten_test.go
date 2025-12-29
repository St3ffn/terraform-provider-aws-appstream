// Copyright (c) St3ffn
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
			prior: mustSet(t, storageConnectorObjectType, []storageConnectorModel{
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
			prior: mustSet(t, storageConnectorObjectType, []storageConnectorModel{
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

			var models []storageConnectorModel
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
	ctx := context.Background()

	prior := mustSet(t, userSettingObjectType, []userSettingModel{
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

	var models []userSettingModel
	diags = out.ElementsAs(ctx, &models, false)
	require.False(t, diags.HasError())
	require.Len(t, models, 1)

	require.Equal(t, "CLIPBOARD_COPY_FROM_LOCAL_DEVICE", models[0].Action.ValueString())
	require.Equal(t, "ENABLED", models[0].Permission.ValueString())
	require.True(t, models[0].MaximumLength.IsNull())
}

func TestFlattenApplicationSettingsResource_computed_bucket_preserved(t *testing.T) {
	ctx := context.Background()

	prior := mustObject(t, applicationSettingsObjectType.AttrTypes, applicationSettingsModel{
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

	var model applicationSettingsModel
	diags = out.As(ctx, &model, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError())

	require.False(t, model.Enabled.ValueBool())
	require.True(t, model.SettingsGroup.IsNull())
	require.Equal(t, "bucket", model.S3BucketName.ValueString())
}

func TestFlattenAccessEndpointsResource_basic_mapping(t *testing.T) {
	ctx := context.Background()

	prior := mustSet(t, accessEndpointObjectType, []accessEndpointModel{
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

	var models []accessEndpointModel
	diags = out.ElementsAs(ctx, &models, false)
	require.False(t, diags.HasError())
	require.Len(t, models, 1)

	require.Equal(t, "STREAMING", models[0].EndpointType.ValueString())
	require.Equal(t, "vpce-123", models[0].VpceID.ValueString())
}

func TestFlattenStreamingExperienceSettingsResource_preferred_protocol_from_aws(t *testing.T) {
	ctx := context.Background()

	prior := mustObject(
		t,
		streamingExperienceSettingsObjectType.AttrTypes,
		streamingExperienceSettingsModel{
			PreferredProtocol: types.StringValue("TCP"),
		},
	)

	awsSettings := &awstypes.StreamingExperienceSettings{
		PreferredProtocol: awstypes.PreferredProtocolUdp,
	}

	var diags diag.Diagnostics
	out := flattenStreamingExperienceSettingsResource(ctx, prior, awsSettings, &diags)
	require.False(t, diags.HasError())

	var model streamingExperienceSettingsModel
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
