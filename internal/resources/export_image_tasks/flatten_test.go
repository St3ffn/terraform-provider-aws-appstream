// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package export_image_tasks

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

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

func TestFlattenExportImageTasksData(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name           string
		aws            []awstypes.ExportImageTask
		expectNull     bool
		expectElements int
	}{
		{
			name:       "empty_returns_null",
			expectNull: true,
		},
		{
			name: "basic_mapping",
			aws: []awstypes.ExportImageTask{
				{
					TaskId:      aws.String("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
					ImageArn:    aws.String("arn:aws:appstream:eu-central-1:123456789012:image/example"),
					AmiName:     aws.String("example-ami"),
					AmiId:       aws.String("ami-0123456789abcdef0"),
					State:       awstypes.ExportImageTaskStateCompleted,
					CreatedDate: &now,
					TagSpecifications: map[string]string{
						"Environment": "test",
					},
				},
			},
			expectElements: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			out := flattenExportImageTasksData(ctx, tt.aws, &diags)
			require.False(t, diags.HasError(), "unexpected diagnostics")

			if tt.expectNull {
				require.True(t, out.IsNull())
				return
			}

			var models []exportImageTaskModel
			diags = out.ElementsAs(ctx, &models, false)
			require.False(t, diags.HasError())
			require.Len(t, models, tt.expectElements)

			task := models[0]
			require.Equal(t, "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", task.TaskID.ValueString())
			require.Equal(t, "example-ami", task.AmiName.ValueString())
			require.Equal(t, "COMPLETED", task.State.ValueString())
			require.Equal(t, "ami-0123456789abcdef0", task.AmiID.ValueString())
		})
	}
}
