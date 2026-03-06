// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package imported_image

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

func TestExpandRuntimeValidationConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		obj  types.Object
		want *awstypes.RuntimeValidationConfig
	}{
		{
			name: "instance_type_set",
			obj: types.ObjectValueMust(
				map[string]attr.Type{
					"intended_instance_type": types.StringType,
				},
				map[string]attr.Value{
					"intended_instance_type": types.StringValue("stream.standard.large"),
				},
			),
			want: &awstypes.RuntimeValidationConfig{
				IntendedInstanceType: aws.String("stream.standard.large"),
			},
		},
		{
			name: "instance_type_null",
			obj: types.ObjectValueMust(
				map[string]attr.Type{
					"intended_instance_type": types.StringType,
				},
				map[string]attr.Value{
					"intended_instance_type": types.StringNull(),
				},
			),
			want: &awstypes.RuntimeValidationConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := expandRuntimeValidationConfig(ctx, tt.obj, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExpandAppCatalogConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	appCatalogConfigElemType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":                   types.StringType,
			"absolute_app_path":      types.StringType,
			"display_name":           types.StringType,
			"launch_parameters":      types.StringType,
			"working_directory":      types.StringType,
			"absolute_icon_path":     types.StringType,
			"absolute_manifest_path": types.StringType,
		},
	}

	tests := []struct {
		name string
		set  types.Set
		want []awstypes.ApplicationConfig
	}{
		{
			name: "single_config",
			set: types.SetValueMust(
				appCatalogConfigElemType,
				[]attr.Value{
					types.ObjectValueMust(
						appCatalogConfigElemType.AttrTypes,
						map[string]attr.Value{
							"name":                   types.StringValue("calculator"),
							"absolute_app_path":      types.StringValue("C:\\Program Files\\App\\app.exe"),
							"display_name":           types.StringValue("Calculator"),
							"launch_parameters":      types.StringValue("--safe-mode"),
							"working_directory":      types.StringValue("C:\\Program Files\\App"),
							"absolute_icon_path":     types.StringValue("C:\\icons\\app.ico"),
							"absolute_manifest_path": types.StringValue("C:\\manifests\\app.json"),
						},
					),
				},
			),
			want: []awstypes.ApplicationConfig{
				{
					Name:                 aws.String("calculator"),
					AbsoluteAppPath:      aws.String("C:\\Program Files\\App\\app.exe"),
					DisplayName:          aws.String("Calculator"),
					LaunchParameters:     aws.String("--safe-mode"),
					WorkingDirectory:     aws.String("C:\\Program Files\\App"),
					AbsoluteIconPath:     aws.String("C:\\icons\\app.ico"),
					AbsoluteManifestPath: aws.String("C:\\manifests\\app.json"),
				},
			},
		},
		{
			name: "empty_set_returns_nil",
			set:  types.SetNull(appCatalogConfigElemType),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := expandAppCatalogConfig(ctx, tt.set, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
