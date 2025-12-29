// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

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

func TestExpandIconS3Location(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     types.Object
		want      *awstypes.S3Location
		wantError bool
	}{
		{
			name: "valid_bucket_and_key",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"s3_bucket": types.StringType,
					"s3_key":    types.StringType,
				},
				map[string]attr.Value{
					"s3_bucket": types.StringValue("my-bucket"),
					"s3_key":    types.StringValue("icon.png"),
				},
			),
			want: &awstypes.S3Location{
				S3Bucket: aws.String("my-bucket"),
				S3Key:    aws.String("icon.png"),
			},
		},
		{
			name: "null_fields",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"s3_bucket": types.StringType,
					"s3_key":    types.StringType,
				},
				map[string]attr.Value{
					"s3_bucket": types.StringNull(),
					"s3_key":    types.StringNull(),
				},
			),
			want: &awstypes.S3Location{
				S3Bucket: nil,
				S3Key:    nil,
			},
		},
		{
			name:      "unknown_object",
			input:     types.ObjectUnknown(map[string]attr.Type{}),
			want:      nil,
			wantError: true,
		},
		{
			name: "invalid_object_shape",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"wrong": types.StringType,
				},
				map[string]attr.Value{
					"wrong": types.StringValue("oops"),
				},
			),
			want:      nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandIconS3Location(ctx, tt.input, &diags)

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
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExpandPlatforms(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     types.Set
		want      []awstypes.PlatformType
		wantError bool
	}{
		{
			name:  "null_set",
			input: types.SetNull(types.StringType),
			want:  nil,
		},
		{
			name:  "unknown_set",
			input: types.SetUnknown(types.StringType),
			want:  nil,
		},
		{
			name:  "empty_set",
			input: types.SetValueMust(types.StringType, []attr.Value{}),
			want:  nil,
		},
		{
			name: "single_platform",
			input: types.SetValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("WINDOWS")},
			),
			want: []awstypes.PlatformType{
				awstypes.PlatformType("WINDOWS"),
			},
		},
		{
			name: "multiple_platforms",
			input: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("WINDOWS"),
					types.StringValue("LINUX"),
				},
			),
			want: []awstypes.PlatformType{
				awstypes.PlatformType("WINDOWS"),
				awstypes.PlatformType("LINUX"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandPlatforms(ctx, tt.input, &diags)

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
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
