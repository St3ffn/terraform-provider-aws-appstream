// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package export_image_task

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenAppBlockErrorsData(t *testing.T) {
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
			want:  types.SetNull(errorDetailObjectType),
		},
		{
			name:  "empty_slice_returns_null_set",
			input: []awstypes.ErrorDetails{},
			want:  types.SetNull(errorDetailObjectType),
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
				errorDetailObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						errorDetailObjectType.AttrTypes,
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
				errorDetailObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						errorDetailObjectType.AttrTypes,
						map[string]attr.Value{
							"error_code":    types.StringValue("ErrorA"),
							"error_message": types.StringValue("Message A"),
						},
					),
					types.ObjectValueMust(
						errorDetailObjectType.AttrTypes,
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
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenErrorDetailsData(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics for test %q: %v", tt.name, diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("flattenErrorDetailsData mismatch for test %q\n got:  %v\n want: %v",
					tt.name, got, tt.want,
				)
			}
		})
	}
}
