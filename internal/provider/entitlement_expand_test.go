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
)

func TestExpandEntitlementAttributes(t *testing.T) {
	ctx := context.Background()

	validSet := types.SetValueMust(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			},
		},
		[]attr.Value{
			types.ObjectValueMust(
				map[string]attr.Type{
					"name":  types.StringType,
					"value": types.StringType,
				},
				map[string]attr.Value{
					"name":  types.StringValue("key1"),
					"value": types.StringValue("val1"),
				},
			),
		},
	)

	tests := []struct {
		name      string
		input     types.Set
		want      []awstypes.EntitlementAttribute
		wantError bool
	}{
		{
			name:  "valid_attributes",
			input: validSet,
			want: []awstypes.EntitlementAttribute{
				{
					Name:  aws.String("key1"),
					Value: aws.String("val1"),
				},
			},
		},
		{
			name: "null_name_returns_error",
			input: types.SetValueMust(
				validSet.ElementType(ctx),
				[]attr.Value{
					types.ObjectValueMust(
						map[string]attr.Type{
							"name":  types.StringType,
							"value": types.StringType,
						},
						map[string]attr.Value{
							"name":  types.StringNull(),
							"value": types.StringValue("val"),
						},
					),
				},
			),
			wantError: true,
		},
		{
			name: "unknown_value_returns_error",
			input: types.SetValueMust(
				validSet.ElementType(ctx),
				[]attr.Value{
					types.ObjectValueMust(
						map[string]attr.Type{
							"name":  types.StringType,
							"value": types.StringType,
						},
						map[string]attr.Value{
							"name":  types.StringValue("key"),
							"value": types.StringUnknown(),
						},
					),
				},
			),
			wantError: true,
		},
		{
			name:  "empty_set_returns_empty_slice",
			input: types.SetValueMust(validSet.ElementType(ctx), []attr.Value{}),
			want:  []awstypes.EntitlementAttribute{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandEntitlementAttributes(ctx, tt.input, &diags)

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
				t.Fatalf("expandEntitlementAttributes() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
