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

func TestFlattenScriptDetailsResource(t *testing.T) {
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

func testScriptDetailsObject(t *testing.T, v map[string]attr.Value) types.Object {
	t.Helper()
	return types.ObjectValueMust(scriptDetailsObjectType.AttrTypes, v)
}
