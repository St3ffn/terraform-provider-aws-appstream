// Copyright St3ffn 2025, 2026
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

func TestFlattenSourceS3LocationResource(t *testing.T) {
	t.Parallel()

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenSourceS3LocationResource(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics for test %q: %v", tt.name, diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("flattenSourceS3LocationResource mismatch for test %q\n got:  %v\n want: %v",
					tt.name, got, tt.want,
				)
			}
		})
	}
}

func TestFlattenScriptDetailsResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	fullAWS := &awstypes.ScriptDetails{
		ExecutablePath:       aws.String("C:\\setup.ps1"),
		ExecutableParameters: aws.String("-flag"),
		TimeoutInSeconds:     aws.Int32(300),
		ScriptS3Location: &awstypes.S3Location{
			S3Bucket: aws.String("bucket"),
			S3Key:    aws.String("script.ps1"),
		},
	}

	expectedFull := testScriptDetailsObject(t, map[string]attr.Value{
		"executable_path":       types.StringValue("C:\\setup.ps1"),
		"executable_parameters": types.StringValue("-flag"),
		"timeout_in_seconds":    types.Int32Value(300),
		"script_s3_location": types.ObjectValueMust(
			s3LocationObjectType.AttrTypes,
			map[string]attr.Value{
				"s3_bucket": types.StringValue("bucket"),
				"s3_key":    types.StringValue("script.ps1"),
			},
		),
	})

	tests := []struct {
		name  string
		prior types.Object
		aws   *awstypes.ScriptDetails
		want  types.Object
	}{
		{
			name:  "prior_null_ignores_aws",
			prior: types.ObjectNull(scriptDetailsObjectType.AttrTypes),
			aws:   fullAWS,
			want:  types.ObjectNull(scriptDetailsObjectType.AttrTypes),
		},
		{
			name:  "prior_unknown_preserved",
			prior: types.ObjectUnknown(scriptDetailsObjectType.AttrTypes),
			aws:   fullAWS,
			want:  types.ObjectUnknown(scriptDetailsObjectType.AttrTypes),
		},
		{
			name: "owned_but_aws_nil_returns_null",
			prior: testScriptDetailsObject(t, map[string]attr.Value{
				"executable_path":       types.StringValue("x"),
				"executable_parameters": types.StringNull(),
				"timeout_in_seconds":    types.Int32Value(10),
				"script_s3_location":    types.ObjectNull(s3LocationObjectType.AttrTypes),
			}),
			aws:  nil,
			want: types.ObjectNull(scriptDetailsObjectType.AttrTypes),
		},
		{
			name: "owned_and_aws_present_flattens",
			prior: testScriptDetailsObject(t, map[string]attr.Value{
				"executable_path":       types.StringValue("old"),
				"executable_parameters": types.StringValue("old"),
				"timeout_in_seconds":    types.Int32Value(1),
				"script_s3_location":    types.ObjectNull(s3LocationObjectType.AttrTypes),
			}),
			aws:  fullAWS,
			want: expectedFull,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenScriptDetailsResource(ctx, tt.prior, tt.aws, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf(
					"flattenScriptDetailsResource mismatch\nprior: %v\naws: %#v\ngot:  %v\nwant: %v",
					tt.prior, tt.aws, got, tt.want,
				)
			}
		})
	}
}

func TestFlattenAppBlockErrorsResource(t *testing.T) {
	t.Parallel()

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenAppBlockErrorsResource(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics for test %q: %v", tt.name, diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("flattenAppBlockErrorsResource mismatch for test %q\n got:  %v\n want: %v",
					tt.name, got, tt.want,
				)
			}
		})
	}
}

func testScriptDetailsObject(t *testing.T, v map[string]attr.Value) types.Object {
	t.Helper()
	return types.ObjectValueMust(scriptDetailsObjectType.AttrTypes, v)
}
