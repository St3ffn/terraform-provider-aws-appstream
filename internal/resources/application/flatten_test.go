// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenIconS3Location(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input *awstypes.S3Location
		want  types.Object
	}{
		{
			name:  "nil_location_returns_null_object",
			input: nil,
			want:  types.ObjectNull(iconS3LocationObjectType.AttrTypes),
		},
		{
			name: "bucket_and_key_set",
			input: &awstypes.S3Location{
				S3Bucket: aws.String("my-bucket"),
				S3Key:    aws.String("icon.png"),
			},
			want: types.ObjectValueMust(
				iconS3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringValue("my-bucket"),
					"s3_key":    types.StringValue("icon.png"),
				},
			),
		},
		{
			name: "partial_values",
			input: &awstypes.S3Location{
				S3Bucket: nil,
				S3Key:    aws.String("icon.png"),
			},
			want: types.ObjectValueMust(
				iconS3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringNull(),
					"s3_key":    types.StringValue("icon.png"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenIconS3Location(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenPlatforms(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input []awstypes.PlatformType
		want  types.Set
	}{
		{
			name:  "nil_slice_returns_empty_set",
			input: nil,
			want:  types.SetValueMust(types.StringType, []attr.Value{}),
		},
		{
			name:  "empty_slice_returns_empty_set",
			input: []awstypes.PlatformType{},
			want:  types.SetValueMust(types.StringType, []attr.Value{}),
		},
		{
			name: "single_platform",
			input: []awstypes.PlatformType{
				awstypes.PlatformType("WINDOWS"),
			},
			want: types.SetValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("WINDOWS")},
			),
		},
		{
			name: "multiple_platforms",
			input: []awstypes.PlatformType{
				awstypes.PlatformType("WINDOWS"),
				awstypes.PlatformType("LINUX"),
			},
			want: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("WINDOWS"),
					types.StringValue("LINUX"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenPlatforms(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
