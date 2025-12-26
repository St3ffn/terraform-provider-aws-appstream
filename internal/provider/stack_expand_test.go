// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestExpandStorageConnectors(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input types.Set
		want  []awstypes.StorageConnector
	}{
		{
			name:  "empty_set",
			input: types.SetValueMust(stackStorageConnectorObjectType, []attr.Value{}),
			want:  []awstypes.StorageConnector{},
		},
		{
			name: "single_connector_with_all_fields",
			input: types.SetValueMust(
				stackStorageConnectorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						stackStorageConnectorObjectType.AttrTypes,
						map[string]attr.Value{
							"connector_type":      types.StringValue("HOMEFOLDERS"),
							"resource_identifier": types.StringValue("arn:aws:s3:::bucket"),
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
			want: []awstypes.StorageConnector{
				{
					ConnectorType:      awstypes.StorageConnectorTypeHomefolders,
					ResourceIdentifier: aws.String("arn:aws:s3:::bucket"),
					Domains:            []string{"example.com"},
					DomainsRequireAdminConsent: []string{
						"admin.example.com",
					},
				},
			},
		},
		{
			name: "optional_fields_null",
			input: types.SetValueMust(
				stackStorageConnectorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						stackStorageConnectorObjectType.AttrTypes,
						map[string]attr.Value{
							"connector_type":                types.StringValue("HOMEFOLDERS"),
							"resource_identifier":           types.StringNull(),
							"domains":                       types.SetNull(types.StringType),
							"domains_require_admin_consent": types.SetNull(types.StringType),
						},
					),
				},
			),
			want: []awstypes.StorageConnector{
				{
					ConnectorType: awstypes.StorageConnectorTypeHomefolders,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandStorageConnectors(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf(
					"expandStorageConnectors() = %#v, want %#v",
					got, tt.want,
				)
			}
		})
	}
}

func TestExpandUserSettings(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     types.Set
		want      []awstypes.UserSetting
		wantError bool
	}{
		{
			name: "empty_set",
			input: types.SetValueMust(
				stackUserSettingObjectType,
				[]attr.Value{},
			),
			want: []awstypes.UserSetting{},
		},
		{
			name: "single_setting_without_max_length",
			input: types.SetValueMust(
				stackUserSettingObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						stackUserSettingObjectType.AttrTypes,
						map[string]attr.Value{
							"action":         types.StringValue("CLIPBOARD_COPY_FROM_LOCAL_DEVICE"),
							"permission":     types.StringValue("ENABLED"),
							"maximum_length": types.Int32Null(),
						},
					),
				},
			),
			want: []awstypes.UserSetting{
				{
					Action:     awstypes.Action("CLIPBOARD_COPY_FROM_LOCAL_DEVICE"),
					Permission: awstypes.Permission("ENABLED"),
				},
			},
		},
		{
			name: "single_setting_with_max_length",
			input: types.SetValueMust(
				stackUserSettingObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						stackUserSettingObjectType.AttrTypes,
						map[string]attr.Value{
							"action":         types.StringValue("CLIPBOARD_COPY_TO_LOCAL_DEVICE"),
							"permission":     types.StringValue("DISABLED"),
							"maximum_length": types.Int32Value(100),
						},
					),
				},
			),
			want: []awstypes.UserSetting{
				{
					Action:        awstypes.Action("CLIPBOARD_COPY_TO_LOCAL_DEVICE"),
					Permission:    awstypes.Permission("DISABLED"),
					MaximumLength: aws.Int32(100),
				},
			},
		},
		{
			name: "invalid_set_element_type",
			input: types.SetValueMust(
				types.StringType, // wrong element type
				[]attr.Value{
					types.StringValue("invalid"),
				},
			),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandUserSettings(ctx, tt.input, &diags)

			if tt.wantError {
				if !diags.HasError() {
					t.Fatalf("expected diagnostics error, got none")
				}
				return
			}

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expandUserSettings() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExpandApplicationSettings(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     types.Object
		want      *awstypes.ApplicationSettings
		wantError bool
	}{
		{
			name: "enabled_with_settings_group",
			input: types.ObjectValueMust(
				stackApplicationSettingsObjectType.AttrTypes,
				map[string]attr.Value{
					"enabled":        types.BoolValue(true),
					"settings_group": types.StringValue("group1"),
					"s3_bucket_name": types.StringNull(), // computed, ignored
				},
			),
			want: &awstypes.ApplicationSettings{
				Enabled:       aws.Bool(true),
				SettingsGroup: aws.String("group1"),
			},
		},
		{
			name: "enabled_without_settings_group",
			input: types.ObjectValueMust(
				stackApplicationSettingsObjectType.AttrTypes,
				map[string]attr.Value{
					"enabled":        types.BoolValue(false),
					"settings_group": types.StringNull(),
					"s3_bucket_name": types.StringNull(),
				},
			),
			want: &awstypes.ApplicationSettings{
				Enabled: aws.Bool(false),
			},
		},
		{
			name: "unknown_enabled",
			input: types.ObjectValueMust(
				stackApplicationSettingsObjectType.AttrTypes,
				map[string]attr.Value{
					"enabled":        types.BoolUnknown(),
					"settings_group": types.StringValue("group1"),
					"s3_bucket_name": types.StringNull(),
				},
			),
			want: &awstypes.ApplicationSettings{
				SettingsGroup: aws.String("group1"),
			},
		},
		{
			name: "invalid_object_type",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"enabled": types.BoolType,
				},
				map[string]attr.Value{
					"enabled": types.BoolValue(true),
				},
			),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandApplicationSettings(ctx, tt.input, &diags)

			if tt.wantError {
				require.True(t, diags.HasError(), "expected diagnostics error")
				return
			}

			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestExpandAccessEndpoints(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     types.Set
		want      []awstypes.AccessEndpoint
		wantError bool
	}{
		{
			name: "empty_set",
			input: types.SetValueMust(
				stackAccessEndpointObjectType,
				[]attr.Value{},
			),
			want: []awstypes.AccessEndpoint{},
		},
		{
			name: "single_endpoint_with_vpce",
			input: types.SetValueMust(
				stackAccessEndpointObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						stackAccessEndpointObjectType.AttrTypes,
						map[string]attr.Value{
							"endpoint_type": types.StringValue("STREAMING"),
							"vpce_id":       types.StringValue("vpce-123"),
						},
					),
				},
			),
			want: []awstypes.AccessEndpoint{
				{
					EndpointType: awstypes.AccessEndpointType("STREAMING"),
					VpceId:       aws.String("vpce-123"),
				},
			},
		},
		{
			name: "endpoint_without_vpce",
			input: types.SetValueMust(
				stackAccessEndpointObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						stackAccessEndpointObjectType.AttrTypes,
						map[string]attr.Value{
							"endpoint_type": types.StringValue("STREAMING"),
							"vpce_id":       types.StringNull(),
						},
					),
				},
			),
			want: []awstypes.AccessEndpoint{
				{
					EndpointType: awstypes.AccessEndpointType("STREAMING"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandAccessEndpoints(ctx, tt.input, &diags)

			if tt.wantError {
				if !diags.HasError() {
					t.Fatalf("expected diagnostics error, got none")
				}
				return
			}

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expandAccessEndpoints() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExpandStreamingExperienceSettings(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     types.Object
		want      *awstypes.StreamingExperienceSettings
		wantError bool
	}{
		{
			name: "null_object_returns_error",
			input: types.ObjectNull(
				stackStreamingExperienceSettingsObjectType.AttrTypes,
			),
			want:      nil,
			wantError: true,
		},
		{
			name: "preferred_protocol_null_returns_nil",
			input: types.ObjectValueMust(
				stackStreamingExperienceSettingsObjectType.AttrTypes,
				map[string]attr.Value{
					"preferred_protocol": types.StringNull(),
				},
			),
			want: nil,
		},
		{
			name: "preferred_protocol_set_returns_struct",
			input: types.ObjectValueMust(
				stackStreamingExperienceSettingsObjectType.AttrTypes,
				map[string]attr.Value{
					"preferred_protocol": types.StringValue("UDP"),
				},
			),
			want: &awstypes.StreamingExperienceSettings{
				PreferredProtocol: awstypes.PreferredProtocol("UDP"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandStreamingExperienceSettings(ctx, tt.input, &diags)

			if tt.wantError {
				if !diags.HasError() {
					t.Fatalf("expected diagnostics error, got none")
				}
				return
			}

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %#v", got)
				}
				return
			}

			if got == nil {
				t.Fatalf("expected %#v, got nil", tt.want)
			}

			if got.PreferredProtocol != tt.want.PreferredProtocol {
				t.Fatalf(
					"PreferredProtocol = %q, want %q",
					got.PreferredProtocol,
					tt.want.PreferredProtocol,
				)
			}
		})
	}
}
