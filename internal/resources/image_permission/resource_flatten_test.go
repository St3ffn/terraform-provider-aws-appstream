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

func TestFlattenImagePermissionsResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.ImagePermissions
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(imagePermissionsObjectType.AttrTypes),
		},
		{
			name: "values_set_both_true",
			in: &awstypes.ImagePermissions{
				AllowFleet:        aws.Bool(true),
				AllowImageBuilder: aws.Bool(true),
			},
			want: types.ObjectValueMust(
				imagePermissionsObjectType.AttrTypes,
				map[string]attr.Value{
					"allow_fleet":         types.BoolValue(true),
					"allow_image_builder": types.BoolValue(true),
				},
			),
		},
		{
			name: "values_set_partial",
			in: &awstypes.ImagePermissions{
				AllowFleet: aws.Bool(false),
				// AllowImageBuilder intentionally nil
			},
			want: types.ObjectValueMust(
				imagePermissionsObjectType.AttrTypes,
				map[string]attr.Value{
					"allow_fleet":         types.BoolValue(false),
					"allow_image_builder": types.BoolNull(),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenImagePermissionsResource(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
