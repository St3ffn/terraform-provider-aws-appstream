// Copyright St3ffn 2025, 2026
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

func TestFlattenImagePermissionsData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		in   []awstypes.SharedImagePermissions
		want types.Set
	}{
		{
			name: "empty_input",
			in:   nil,
			want: types.SetNull(imagePermissionEntryObjectType),
		},
		{
			name: "single_entry",
			in: []awstypes.SharedImagePermissions{
				{
					SharedAccountId: aws.String("123456789012"),
					ImagePermissions: &awstypes.ImagePermissions{
						AllowFleet:        aws.Bool(true),
						AllowImageBuilder: aws.Bool(false),
					},
				},
			},
			want: types.SetValueMust(
				imagePermissionEntryObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						imagePermissionEntryObjectType.AttrTypes,
						map[string]attr.Value{
							"shared_account_id": types.StringValue("123456789012"),
							"image_permissions": types.ObjectValueMust(
								imagePermissionsObjectType.AttrTypes,
								map[string]attr.Value{
									"allow_fleet":         types.BoolValue(true),
									"allow_image_builder": types.BoolValue(false),
								},
							),
						},
					),
				},
			),
		},
		{
			name: "multiple_entries",
			in: []awstypes.SharedImagePermissions{
				{
					SharedAccountId: aws.String("111111111111"),
					ImagePermissions: &awstypes.ImagePermissions{
						AllowFleet: aws.Bool(true),
					},
				},
				{
					SharedAccountId: aws.String("222222222222"),
					ImagePermissions: &awstypes.ImagePermissions{
						AllowImageBuilder: aws.Bool(true),
					},
				},
			},
			want: types.SetValueMust(
				imagePermissionEntryObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						imagePermissionEntryObjectType.AttrTypes,
						map[string]attr.Value{
							"shared_account_id": types.StringValue("111111111111"),
							"image_permissions": types.ObjectValueMust(
								imagePermissionsObjectType.AttrTypes,
								map[string]attr.Value{
									"allow_fleet":         types.BoolValue(true),
									"allow_image_builder": types.BoolNull(),
								},
							),
						},
					),
					types.ObjectValueMust(
						imagePermissionEntryObjectType.AttrTypes,
						map[string]attr.Value{
							"shared_account_id": types.StringValue("222222222222"),
							"image_permissions": types.ObjectValueMust(
								imagePermissionsObjectType.AttrTypes,
								map[string]attr.Value{
									"allow_fleet":         types.BoolNull(),
									"allow_image_builder": types.BoolValue(true),
								},
							),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenImagePermissionsData(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
