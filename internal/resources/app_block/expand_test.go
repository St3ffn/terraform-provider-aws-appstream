// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

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

func TestExpandSourceS3Location(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     types.Object
		want      *awstypes.S3Location
		wantError bool
	}{
		{
			name: "bucket_and_key_set",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"s3_bucket": types.StringType,
					"s3_key":    types.StringType,
				},
				map[string]attr.Value{
					"s3_bucket": types.StringValue("my-bucket"),
					"s3_key":    types.StringValue("path/to/object"),
				},
			),
			want: &awstypes.S3Location{
				S3Bucket: aws.String("my-bucket"),
				S3Key:    aws.String("path/to/object"),
			},
		},
		{
			name: "bucket_only",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"s3_bucket": types.StringType,
					"s3_key":    types.StringType,
				},
				map[string]attr.Value{
					"s3_bucket": types.StringValue("my-bucket"),
					"s3_key":    types.StringNull(),
				},
			),
			want: &awstypes.S3Location{
				S3Bucket: aws.String("my-bucket"),
				S3Key:    nil,
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

			got := expandSourceS3Location(ctx, tt.input, &diags)

			if tt.wantError {
				if !diags.HasError() {
					t.Fatalf("expected diagnostics error for test %q, got none", tt.name)
				}
				return
			}

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics for test %q: %v", tt.name, diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("mismatch for test %q:\n got:  %#v\n want: %#v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpandScriptDetails(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     types.Object
		want      *awstypes.ScriptDetails
		wantError bool
	}{
		{
			name: "full_script_details",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"script_s3_location":    s3LocationObjectType,
					"executable_path":       types.StringType,
					"executable_parameters": types.StringType,
					"timeout_in_seconds":    types.Int32Type,
				},
				map[string]attr.Value{
					"script_s3_location": types.ObjectValueMust(
						map[string]attr.Type{
							"s3_bucket": types.StringType,
							"s3_key":    types.StringType,
						},
						map[string]attr.Value{
							"s3_bucket": types.StringValue("script-bucket"),
							"s3_key":    types.StringValue("setup.ps1"),
						},
					),
					"executable_path":       types.StringValue("C:\\setup.ps1"),
					"executable_parameters": types.StringValue("-flag"),
					"timeout_in_seconds":    types.Int32Value(300),
				},
			),
			want: &awstypes.ScriptDetails{
				ExecutablePath:       aws.String("C:\\setup.ps1"),
				ExecutableParameters: aws.String("-flag"),
				TimeoutInSeconds:     aws.Int32(300),
				ScriptS3Location: &awstypes.S3Location{
					S3Bucket: aws.String("script-bucket"),
					S3Key:    aws.String("setup.ps1"),
				},
			},
		},
		{
			name: "script_without_s3_location",
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"script_s3_location":    s3LocationObjectType,
					"executable_path":       types.StringType,
					"executable_parameters": types.StringType,
					"timeout_in_seconds":    types.Int32Type,
				},
				map[string]attr.Value{
					"script_s3_location":    types.ObjectNull(s3LocationObjectType.AttrTypes),
					"executable_path":       types.StringValue("/bin/run.sh"),
					"executable_parameters": types.StringNull(),
					"timeout_in_seconds":    types.Int32Value(60),
				},
			),
			want: &awstypes.ScriptDetails{
				ExecutablePath:       aws.String("/bin/run.sh"),
				ExecutableParameters: nil,
				TimeoutInSeconds:     aws.Int32(60),
				ScriptS3Location:     nil,
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

			got := expandScriptDetails(ctx, tt.input, &diags)

			if tt.wantError {
				if !diags.HasError() {
					t.Fatalf("expected diagnostics error for test %q, got none", tt.name)
				}
				return
			}

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics for test %q: %v", tt.name, diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("mismatch for test %q:\n got:  %#v\n want: %#v", tt.name, got, tt.want)
			}
		})
	}
}
