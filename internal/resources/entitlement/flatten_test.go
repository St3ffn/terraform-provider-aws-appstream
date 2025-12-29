// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenAttributes(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input []awstypes.EntitlementAttribute
		want  types.Set
	}{
		{
			name: "multiple_attributes",
			input: []awstypes.EntitlementAttribute{
				{
					Name:  aws.String("key1"),
					Value: aws.String("val1"),
				},
				{
					Name:  aws.String("key2"),
					Value: aws.String("val2"),
				},
			},
			want: types.SetValueMust(
				attributeObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						attributeObjectType.AttrTypes,
						map[string]attr.Value{
							"name":  types.StringValue("key1"),
							"value": types.StringValue("val1"),
						},
					),
					types.ObjectValueMust(
						attributeObjectType.AttrTypes,
						map[string]attr.Value{
							"name":  types.StringValue("key2"),
							"value": types.StringValue("val2"),
						},
					),
				},
			),
		},
		{
			name:  "empty_slice_returns_empty_set",
			input: []awstypes.EntitlementAttribute{},
			want: types.SetValueMust(
				attributeObjectType,
				[]attr.Value{},
			),
		},
		{
			name:  "nil_slice_returns_empty_set",
			input: nil,
			want: types.SetValueMust(
				attributeObjectType,
				[]attr.Value{},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenAttributes(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf(
					"flattenAttributes() = %#v, want %#v",
					got, tt.want,
				)
			}

			if got.IsNull() {
				t.Fatalf("unexpected null set")
			}
		})
	}
}
