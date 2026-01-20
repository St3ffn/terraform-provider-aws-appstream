// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_permission

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExpandImagePermissions(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   types.Object
		want *awstypes.ImagePermissions
	}{
		{
			name: "values_set_both_true",
			in: types.ObjectValueMust(
				imagePermissionsObjectType.AttrTypes,
				map[string]attr.Value{
					"allow_fleet":         types.BoolValue(true),
					"allow_image_builder": types.BoolValue(true),
				},
			),
			want: &awstypes.ImagePermissions{
				AllowFleet:        aws.Bool(true),
				AllowImageBuilder: aws.Bool(true),
			},
		},
		{
			name: "values_set_both_false",
			in: types.ObjectValueMust(
				imagePermissionsObjectType.AttrTypes,
				map[string]attr.Value{
					"allow_fleet":         types.BoolValue(false),
					"allow_image_builder": types.BoolValue(false),
				},
			),
			want: &awstypes.ImagePermissions{
				AllowFleet:        aws.Bool(false),
				AllowImageBuilder: aws.Bool(false),
			},
		},
		{
			name: "values_null_defensive",
			in: types.ObjectValueMust(
				imagePermissionsObjectType.AttrTypes,
				map[string]attr.Value{
					"allow_fleet":         types.BoolNull(),
					"allow_image_builder": types.BoolNull(),
				},
			),
			want: &awstypes.ImagePermissions{
				AllowFleet:        nil,
				AllowImageBuilder: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandImagePermissions(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if tt.want == nil && got != nil {
				t.Fatalf("got %v, want nil", got)
			}

			if tt.want != nil && got == nil {
				t.Fatalf("got nil, want %v", tt.want)
			}

			if tt.want != nil {
				if (got.AllowFleet == nil) != (tt.want.AllowFleet == nil) {
					t.Fatalf("AllowFleet mismatch: got %v, want %v", got.AllowFleet, tt.want.AllowFleet)
				}
				if got.AllowFleet != nil && *got.AllowFleet != *tt.want.AllowFleet {
					t.Fatalf("AllowFleet value mismatch: got %v, want %v", *got.AllowFleet, *tt.want.AllowFleet)
				}

				if (got.AllowImageBuilder == nil) != (tt.want.AllowImageBuilder == nil) {
					t.Fatalf("AllowImageBuilder mismatch: got %v, want %v", got.AllowImageBuilder, tt.want.AllowImageBuilder)
				}
				if got.AllowImageBuilder != nil && *got.AllowImageBuilder != *tt.want.AllowImageBuilder {
					t.Fatalf("AllowImageBuilder value mismatch: got %v, want %v", *got.AllowImageBuilder, *tt.want.AllowImageBuilder)
				}
			}
		})
	}
}
