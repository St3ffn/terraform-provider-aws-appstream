// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package export_image_tasks

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

func TestExpandFilters(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	filterObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":   types.StringType,
			"values": types.SetType{ElemType: types.StringType},
		},
	}

	validSet := types.SetValueMust(
		filterObjectType,
		[]attr.Value{
			types.ObjectValueMust(
				filterObjectType.AttrTypes,
				map[string]attr.Value{
					"name": types.StringValue("state"),
					"values": types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("COMPLETED"),
							types.StringValue("FAILED"),
						},
					),
				},
			),
		},
	)

	tests := []struct {
		name      string
		input     types.Set
		want      []awstypes.Filter
		wantError bool
	}{
		{
			name:  "valid_filters",
			input: validSet,
			want: []awstypes.Filter{
				{
					Name:   aws.String("state"),
					Values: []string{"COMPLETED", "FAILED"},
				},
			},
		},
		{
			name:  "null_set_returns_nil",
			input: types.SetNull(filterObjectType),
			want:  nil,
		},
		{
			name:  "unknown_set_returns_nil",
			input: types.SetUnknown(filterObjectType),
			want:  nil,
		},
		{
			name: "empty_set_returns_nil",
			input: types.SetValueMust(
				filterObjectType,
				[]attr.Value{},
			),
			want: nil,
		},
		{
			name: "invalid_value_produces_diagnostics",
			input: types.SetValueMust(
				filterObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						filterObjectType.AttrTypes,
						map[string]attr.Value{
							"name": types.StringValue("state"),
							"values": types.SetValueMust(
								types.StringType,
								[]attr.Value{
									types.StringUnknown(),
								},
							),
						},
					),
				},
			),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := expandFilters(ctx, tt.input, &diags)

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
				t.Fatalf("expandFilters() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
