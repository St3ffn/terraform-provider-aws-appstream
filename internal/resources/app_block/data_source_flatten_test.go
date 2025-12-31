// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenSourceS3LocationData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input *awstypes.S3Location
		want  types.Object
	}{
		{
			name:  "nil_location_returns_null_object",
			input: nil,
			want:  types.ObjectNull(sourceS3LocationObjectType.AttrTypes),
		},
		{
			name: "bucket_and_key_set",
			input: &awstypes.S3Location{
				S3Bucket: aws.String("my-bucket"),
				S3Key:    aws.String("path/to/object"),
			},
			want: types.ObjectValueMust(
				sourceS3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringValue("my-bucket"),
					"s3_key":    types.StringValue("path/to/object"),
				},
			),
		},
		{
			name: "bucket_only",
			input: &awstypes.S3Location{
				S3Bucket: aws.String("my-bucket"),
				S3Key:    nil,
			},
			want: types.ObjectValueMust(
				sourceS3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringValue("my-bucket"),
					"s3_key":    types.StringNull(),
				},
			),
		},
		{
			name: "key_only",
			input: &awstypes.S3Location{
				S3Bucket: nil,
				S3Key:    aws.String("object"),
			},
			want: types.ObjectValueMust(
				sourceS3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringNull(),
					"s3_key":    types.StringValue("object"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenSourceS3LocationData(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics for test %q: %v", tt.name, diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("flattenSourceS3Location mismatch for test %q\n got:  %v\n want: %v",
					tt.name, got, tt.want,
				)
			}
		})
	}
}

func TestFlattenScriptDetailsData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input *awstypes.ScriptDetails
		want  types.Object
	}{
		{
			name:  "nil_script_details_returns_null_object",
			input: nil,
			want:  types.ObjectNull(scriptDetailsObjectType.AttrTypes),
		},
		{
			name: "full_script_details",
			input: &awstypes.ScriptDetails{
				ExecutablePath:       aws.String("C:\\setup.ps1"),
				ExecutableParameters: aws.String("-flag"),
				TimeoutInSeconds:     aws.Int32(300),
				ScriptS3Location: &awstypes.S3Location{
					S3Bucket: aws.String("script-bucket"),
					S3Key:    aws.String("setup.ps1"),
				},
			},
			want: types.ObjectValueMust(
				scriptDetailsObjectType.AttrTypes,
				map[string]attr.Value{
					"executable_path":       types.StringValue("C:\\setup.ps1"),
					"executable_parameters": types.StringValue("-flag"),
					"timeout_in_seconds":    types.Int32Value(300),
					"script_s3_location": types.ObjectValueMust(
						s3LocationObjectType.AttrTypes,
						map[string]attr.Value{
							"s3_bucket": types.StringValue("script-bucket"),
							"s3_key":    types.StringValue("setup.ps1"),
						},
					),
				},
			),
		},
		{
			name: "script_without_s3_location",
			input: &awstypes.ScriptDetails{
				ExecutablePath:       aws.String("/bin/run.sh"),
				ExecutableParameters: nil,
				TimeoutInSeconds:     aws.Int32(60),
				ScriptS3Location:     nil,
			},
			want: types.ObjectValueMust(
				scriptDetailsObjectType.AttrTypes,
				map[string]attr.Value{
					"executable_path":       types.StringValue("/bin/run.sh"),
					"executable_parameters": types.StringNull(),
					"timeout_in_seconds":    types.Int32Value(60),
					"script_s3_location":    types.ObjectNull(s3LocationObjectType.AttrTypes),
				},
			),
		},
		{
			name: "script_with_empty_fields",
			input: &awstypes.ScriptDetails{
				ExecutablePath:       nil,
				ExecutableParameters: nil,
				TimeoutInSeconds:     nil,
				ScriptS3Location: &awstypes.S3Location{
					S3Bucket: nil,
					S3Key:    nil,
				},
			},
			want: types.ObjectValueMust(
				scriptDetailsObjectType.AttrTypes,
				map[string]attr.Value{
					"executable_path":       types.StringNull(),
					"executable_parameters": types.StringNull(),
					"timeout_in_seconds":    types.Int32Null(),
					"script_s3_location": types.ObjectValueMust(
						s3LocationObjectType.AttrTypes,
						map[string]attr.Value{
							"s3_bucket": types.StringNull(),
							"s3_key":    types.StringNull(),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenScriptDetailsData(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics for test %q: %v", tt.name, diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("flattenScriptDetails mismatch for test %q\n got:  %v\n want: %v",
					tt.name, got, tt.want,
				)
			}
		})
	}
}

func TestFlattenAppBlockErrorsData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input []awstypes.ErrorDetails
		want  types.Set
	}{
		{
			name:  "nil_slice_returns_null_set",
			input: nil,
			want:  types.SetNull(errorObjectType),
		},
		{
			name:  "empty_slice_returns_null_set",
			input: []awstypes.ErrorDetails{},
			want:  types.SetNull(errorObjectType),
		},
		{
			name: "single_error",
			input: []awstypes.ErrorDetails{
				{
					ErrorCode:    aws.String("InvalidParameter"),
					ErrorMessage: aws.String("Something went wrong"),
				},
			},
			want: types.SetValueMust(
				errorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						errorObjectType.AttrTypes,
						map[string]attr.Value{
							"error_code":    types.StringValue("InvalidParameter"),
							"error_message": types.StringValue("Something went wrong"),
						},
					),
				},
			),
		},
		{
			name: "multiple_errors",
			input: []awstypes.ErrorDetails{
				{
					ErrorCode:    aws.String("ErrorA"),
					ErrorMessage: aws.String("Message A"),
				},
				{
					ErrorCode:    aws.String("ErrorB"),
					ErrorMessage: aws.String("Message B"),
				},
			},
			want: types.SetValueMust(
				errorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						errorObjectType.AttrTypes,
						map[string]attr.Value{
							"error_code":    types.StringValue("ErrorA"),
							"error_message": types.StringValue("Message A"),
						},
					),
					types.ObjectValueMust(
						errorObjectType.AttrTypes,
						map[string]attr.Value{
							"error_code":    types.StringValue("ErrorB"),
							"error_message": types.StringValue("Message B"),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenAppBlockErrorsData(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics for test %q: %v", tt.name, diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("flattenAppBlockErrors mismatch for test %q\n got:  %v\n want: %v",
					tt.name, got, tt.want,
				)
			}
		})
	}
}
