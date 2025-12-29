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
	"github.com/stretchr/testify/require"
)

func TestFlattenStorageConnectorsData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		input   []awstypes.StorageConnector
		want    types.Set
		wantNil bool
	}{
		{
			name:    "empty_input_returns_null_set",
			input:   nil,
			want:    types.SetNull(storageConnectorObjectType),
			wantNil: true,
		},
		{
			name: "single_connector",
			input: []awstypes.StorageConnector{
				{
					ConnectorType:      awstypes.StorageConnectorTypeHomefolders,
					ResourceIdentifier: aws.String("arn:aws:s3:::example"),
					Domains:            []string{"example.com"},
					DomainsRequireAdminConsent: []string{
						"admin.example.com",
					},
				},
			},
			want: types.SetValueMust(
				storageConnectorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						storageConnectorObjectType.AttrTypes,
						map[string]attr.Value{
							"connector_type": types.StringValue("HOMEFOLDERS"),
							"resource_identifier": types.StringValue(
								"arn:aws:s3:::example",
							),
							"domains": types.SetValueMust(
								types.StringType,
								[]attr.Value{
									types.StringValue("example.com"),
								},
							),
							"domains_require_admin_consent": types.SetValueMust(
								types.StringType,
								[]attr.Value{
									types.StringValue("admin.example.com"),
								},
							),
						},
					),
				},
			),
		},
		{
			name: "multiple_connectors",
			input: []awstypes.StorageConnector{
				{
					ConnectorType: awstypes.StorageConnectorTypeHomefolders,
				},
				{
					ConnectorType: awstypes.StorageConnectorTypeGoogleDrive,
				},
			},
			want: types.SetValueMust(
				storageConnectorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						storageConnectorObjectType.AttrTypes,
						map[string]attr.Value{
							"connector_type":                types.StringValue("HOMEFOLDERS"),
							"resource_identifier":           types.StringNull(),
							"domains":                       types.SetNull(types.StringType),
							"domains_require_admin_consent": types.SetNull(types.StringType),
						},
					),
					types.ObjectValueMust(
						storageConnectorObjectType.AttrTypes,
						map[string]attr.Value{
							"connector_type":                types.StringValue("GOOGLE_DRIVE"),
							"resource_identifier":           types.StringNull(),
							"domains":                       types.SetNull(types.StringType),
							"domains_require_admin_consent": types.SetNull(types.StringType),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenStorageConnectorsData(ctx, tt.input, &diags)

			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			if tt.wantNil {
				require.True(t, got.IsNull(), "expected null set")
				return
			}

			require.True(t, got.Equal(tt.want), "flattenStorageConnectors result mismatch")
		})
	}
}

func TestFlattenUserSettingsData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		input   []awstypes.UserSetting
		want    types.Set
		wantNil bool
	}{
		{
			name:    "empty_input_returns_null_set",
			input:   nil,
			want:    types.SetNull(userSettingObjectType),
			wantNil: true,
		},
		{
			name: "single_user_setting_without_max_length",
			input: []awstypes.UserSetting{
				{
					Action:     awstypes.ActionClipboardCopyFromLocalDevice,
					Permission: awstypes.PermissionEnabled,
				},
			},
			want: types.SetValueMust(
				userSettingObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						userSettingObjectType.AttrTypes,
						map[string]attr.Value{
							"action":         types.StringValue("CLIPBOARD_COPY_FROM_LOCAL_DEVICE"),
							"permission":     types.StringValue("ENABLED"),
							"maximum_length": types.Int32Null(),
						},
					),
				},
			),
		},
		{
			name: "single_user_setting_with_max_length",
			input: []awstypes.UserSetting{
				{
					Action:        awstypes.ActionClipboardCopyFromLocalDevice,
					Permission:    awstypes.PermissionEnabled,
					MaximumLength: aws.Int32(100),
				},
			},
			want: types.SetValueMust(
				userSettingObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						userSettingObjectType.AttrTypes,
						map[string]attr.Value{
							"action":         types.StringValue("CLIPBOARD_COPY_FROM_LOCAL_DEVICE"),
							"permission":     types.StringValue("ENABLED"),
							"maximum_length": types.Int32Value(100),
						},
					),
				},
			),
		},
		{
			name: "multiple_user_settings",
			input: []awstypes.UserSetting{
				{
					Action:     awstypes.ActionClipboardCopyFromLocalDevice,
					Permission: awstypes.PermissionEnabled,
				},
				{
					Action:     awstypes.ActionClipboardCopyToLocalDevice,
					Permission: awstypes.PermissionDisabled,
				},
			},
			want: types.SetValueMust(
				userSettingObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						userSettingObjectType.AttrTypes,
						map[string]attr.Value{
							"action":         types.StringValue("CLIPBOARD_COPY_FROM_LOCAL_DEVICE"),
							"permission":     types.StringValue("ENABLED"),
							"maximum_length": types.Int32Null(),
						},
					),
					types.ObjectValueMust(
						userSettingObjectType.AttrTypes,
						map[string]attr.Value{
							"action":         types.StringValue("CLIPBOARD_COPY_TO_LOCAL_DEVICE"),
							"permission":     types.StringValue("DISABLED"),
							"maximum_length": types.Int32Null(),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenUserSettingsData(ctx, tt.input, &diags)

			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			if tt.wantNil {
				require.True(t, got.IsNull(), "expected null set")
				return
			}

			require.True(t, got.Equal(tt.want), "flattenUserSettings result mismatch")
		})
	}
}

func TestFlattenApplicationSettingsData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		input   *awstypes.ApplicationSettingsResponse
		want    types.Object
		wantNil bool
	}{
		{
			name:    "nil_input_returns_null_object",
			input:   nil,
			want:    types.ObjectNull(applicationSettingsObjectType.AttrTypes),
			wantNil: true,
		},
		{
			name: "all_fields_set",
			input: &awstypes.ApplicationSettingsResponse{
				Enabled:       aws.Bool(true),
				SettingsGroup: aws.String("group1"),
				S3BucketName:  aws.String("bucket-name"),
			},
			want: types.ObjectValueMust(
				applicationSettingsObjectType.AttrTypes,
				map[string]attr.Value{
					"enabled":        types.BoolValue(true),
					"settings_group": types.StringValue("group1"),
					"s3_bucket_name": types.StringValue("bucket-name"),
				},
			),
		},
		{
			name: "optional_fields_null",
			input: &awstypes.ApplicationSettingsResponse{
				Enabled: aws.Bool(false),
			},
			want: types.ObjectValueMust(
				applicationSettingsObjectType.AttrTypes,
				map[string]attr.Value{
					"enabled":        types.BoolValue(false),
					"settings_group": types.StringNull(),
					"s3_bucket_name": types.StringNull(),
				},
			),
		},
		{
			name: "enabled_null",
			input: &awstypes.ApplicationSettingsResponse{
				SettingsGroup: aws.String("group1"),
			},
			want: types.ObjectValueMust(
				applicationSettingsObjectType.AttrTypes,
				map[string]attr.Value{
					"enabled":        types.BoolNull(),
					"settings_group": types.StringValue("group1"),
					"s3_bucket_name": types.StringNull(),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenApplicationSettingsData(ctx, tt.input, &diags)

			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			if tt.wantNil {
				require.True(t, got.IsNull(), "expected null object")
				return
			}

			require.True(
				t, got.Equal(tt.want),
				"flattenApplicationSettings() = %#v, want %#v", got, tt.want,
			)
		})
	}
}

func TestFlattenAccessEndpointsData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input []awstypes.AccessEndpoint
		want  types.Set
	}{
		{
			name:  "empty_input_returns_null_set",
			input: nil,
			want:  types.SetNull(accessEndpointObjectType),
		},
		{
			name: "single_endpoint_with_vpce",
			input: []awstypes.AccessEndpoint{
				{
					EndpointType: awstypes.AccessEndpointTypeStreaming,
					VpceId:       aws.String("vpce-123"),
				},
			},
			want: types.SetValueMust(
				accessEndpointObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						accessEndpointObjectType.AttrTypes,
						map[string]attr.Value{
							"endpoint_type": types.StringValue("STREAMING"),
							"vpce_id":       types.StringValue("vpce-123"),
						},
					),
				},
			),
		},
		{
			name: "single_endpoint_without_vpce",
			input: []awstypes.AccessEndpoint{
				{
					EndpointType: awstypes.AccessEndpointTypeStreaming,
				},
			},
			want: types.SetValueMust(
				accessEndpointObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						accessEndpointObjectType.AttrTypes,
						map[string]attr.Value{
							"endpoint_type": types.StringValue("STREAMING"),
							"vpce_id":       types.StringNull(),
						},
					),
				},
			),
		},
		{
			name: "multiple_endpoints",
			input: []awstypes.AccessEndpoint{
				{
					EndpointType: awstypes.AccessEndpointTypeStreaming,
					VpceId:       aws.String("vpce-1"),
				},
				{
					EndpointType: awstypes.AccessEndpointTypeStreaming,
					VpceId:       aws.String("vpce-2"),
				},
			},
			want: types.SetValueMust(
				accessEndpointObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						accessEndpointObjectType.AttrTypes,
						map[string]attr.Value{
							"endpoint_type": types.StringValue("STREAMING"),
							"vpce_id":       types.StringValue("vpce-1"),
						},
					),
					types.ObjectValueMust(
						accessEndpointObjectType.AttrTypes,
						map[string]attr.Value{
							"endpoint_type": types.StringValue("STREAMING"),
							"vpce_id":       types.StringValue("vpce-2"),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenAccessEndpointsData(ctx, tt.input, &diags)

			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			require.True(
				t, got.Equal(tt.want),
				"flattenAccessEndpoints() = %#v, want %#v", got, tt.want,
			)
		})
	}
}

func TestFlattenStreamingExperienceSettingsData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input *awstypes.StreamingExperienceSettings
		want  types.Object
	}{
		{
			name:  "nil_input_returns_null_object",
			input: nil,
			want:  types.ObjectNull(streamingExperienceSettingsObjectType.AttrTypes),
		},
		{
			name: "preferred_protocol_set",
			input: &awstypes.StreamingExperienceSettings{
				PreferredProtocol: awstypes.PreferredProtocolTcp,
			},
			want: types.ObjectValueMust(
				streamingExperienceSettingsObjectType.AttrTypes,
				map[string]attr.Value{
					"preferred_protocol": types.StringValue("TCP"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenStreamingExperienceSettingsData(ctx, tt.input, &diags)

			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			require.True(
				t, got.Equal(tt.want),
				"flattenStreamingExperienceSettings() = %#v, want %#v", got, tt.want,
			)
		})
	}
}

func TestFlattenStackErrorsData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   []awstypes.StackError
		want types.Set
	}{
		{
			name: "empty_slice",
			in:   nil,
			want: types.SetNull(errorObjectType),
		},
		{
			name: "single_error",
			in: []awstypes.StackError{
				{
					ErrorCode:    awstypes.StackErrorCodeInternalServiceError,
					ErrorMessage: aws.String("boom"),
				},
			},
			want: types.SetValueMust(
				errorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						errorObjectType.AttrTypes,
						map[string]attr.Value{
							"error_code":    types.StringValue("INTERNAL_SERVICE_ERROR"),
							"error_message": types.StringValue("boom"),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenStackErrorsData(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
