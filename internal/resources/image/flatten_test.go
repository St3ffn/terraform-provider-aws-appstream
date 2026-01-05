// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenStateChangeReason(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input *awstypes.ImageStateChangeReason
		want  types.Object
	}{
		{
			name:  "nil_reason_returns_null_object",
			input: nil,
			want:  types.ObjectNull(stateChangeReasonObjectType.AttrTypes),
		},
		{
			name: "code_and_message_set",
			input: &awstypes.ImageStateChangeReason{
				Code:    awstypes.ImageStateChangeReasonCodeImageBuilderNotAvailable,
				Message: aws.String("something went wrong"),
			},
			want: types.ObjectValueMust(
				stateChangeReasonObjectType.AttrTypes,
				map[string]attr.Value{
					"code":    types.StringValue(string(awstypes.ImageStateChangeReasonCodeImageBuilderNotAvailable)),
					"message": types.StringValue("something went wrong"),
				},
			),
		},
		{
			name: "code_only",
			input: &awstypes.ImageStateChangeReason{
				Code:    awstypes.ImageStateChangeReasonCodeImageBuilderNotAvailable,
				Message: nil,
			},
			want: types.ObjectValueMust(
				stateChangeReasonObjectType.AttrTypes,
				map[string]attr.Value{
					"code":    types.StringValue(string(awstypes.ImageStateChangeReasonCodeImageBuilderNotAvailable)),
					"message": types.StringNull(),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenStateChangeReason(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenApplications(t *testing.T) {
	ctx := context.Background()

	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		input []awstypes.Application
		want  types.Set
	}{
		{
			name:  "nil_slice_returns_null_set",
			input: nil,
			want:  types.SetNull(applicationObjectType),
		},
		{
			name:  "empty_slice_returns_null_set",
			input: []awstypes.Application{},
			want:  types.SetNull(applicationObjectType),
		},
		{
			name: "single_application",
			input: []awstypes.Application{
				{
					Name:        aws.String("app1"),
					DisplayName: aws.String("App One"),
					Arn:         aws.String("arn:aws:appstream:us-east-1:123:application/app1"),
					Platforms: []awstypes.PlatformType{
						awstypes.PlatformTypeWindows,
					},
					CreatedTime: &ts,
				},
			},
			want: types.SetValueMust(
				applicationObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						applicationObjectType.AttrTypes,
						map[string]attr.Value{
							"name":              types.StringValue("app1"),
							"display_name":      types.StringValue("App One"),
							"icon_url":          types.StringNull(),
							"launch_path":       types.StringNull(),
							"launch_parameters": types.StringNull(),
							"enabled":           types.BoolNull(),
							"metadata":          types.MapNull(types.StringType),
							"working_directory": types.StringNull(),
							"description":       types.StringNull(),
							"arn":               types.StringValue("arn:aws:appstream:us-east-1:123:application/app1"),
							"app_block_arn":     types.StringNull(),
							"icon_s3_location":  types.ObjectNull(iconS3LocationObjectType.AttrTypes),
							"platforms": types.SetValueMust(
								types.StringType,
								[]attr.Value{types.StringValue("WINDOWS")},
							),
							"instance_families": types.SetNull(types.StringType),
							"created_time":      types.StringValue("2024-01-01T00:00:00Z"),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenApplications(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

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
				S3Key:    aws.String("icons/icon.png"),
			},
			want: types.ObjectValueMust(
				iconS3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringValue("my-bucket"),
					"s3_key":    types.StringValue("icons/icon.png"),
				},
			),
		},
		{
			name: "only_bucket_set",
			input: &awstypes.S3Location{
				S3Bucket: aws.String("my-bucket"),
				S3Key:    nil,
			},
			want: types.ObjectValueMust(
				iconS3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringValue("my-bucket"),
					"s3_key":    types.StringNull(),
				},
			),
		},
		{
			name: "only_key_set",
			input: &awstypes.S3Location{
				S3Bucket: nil,
				S3Key:    aws.String("icons/icon.png"),
			},
			want: types.ObjectValueMust(
				iconS3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringNull(),
					"s3_key":    types.StringValue("icons/icon.png"),
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

func TestFlattenImagePermissions(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input *awstypes.ImagePermissions
		want  types.Object
	}{
		{
			name:  "nil_permissions_returns_null_object",
			input: nil,
			want:  types.ObjectNull(imagePermissionsObjectType.AttrTypes),
		},
		{
			name: "both_permissions_set",
			input: &awstypes.ImagePermissions{
				AllowFleet:        aws.Bool(true),
				AllowImageBuilder: aws.Bool(false),
			},
			want: types.ObjectValueMust(
				imagePermissionsObjectType.AttrTypes,
				map[string]attr.Value{
					"allow_fleet":         types.BoolValue(true),
					"allow_image_builder": types.BoolValue(false),
				},
			),
		},
		{
			name: "partial_permissions",
			input: &awstypes.ImagePermissions{
				AllowFleet:        nil,
				AllowImageBuilder: aws.Bool(true),
			},
			want: types.ObjectValueMust(
				imagePermissionsObjectType.AttrTypes,
				map[string]attr.Value{
					"allow_fleet":         types.BoolNull(),
					"allow_image_builder": types.BoolValue(true),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenImagePermissions(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenImageErrors(t *testing.T) {
	ctx := context.Background()

	ts := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name  string
		input []awstypes.ResourceError
		want  types.Set
	}{
		{
			name:  "nil_slice_returns_null_set",
			input: nil,
			want:  types.SetNull(imageErrorObjectType),
		},
		{
			name:  "empty_slice_returns_null_set",
			input: []awstypes.ResourceError{},
			want:  types.SetNull(imageErrorObjectType),
		},
		{
			name: "single_error",
			input: []awstypes.ResourceError{
				{
					ErrorCode:      awstypes.FleetErrorCodeImageNotFound,
					ErrorMessage:   aws.String("image missing"),
					ErrorTimestamp: &ts,
				},
			},
			want: types.SetValueMust(
				imageErrorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						imageErrorObjectType.AttrTypes,
						map[string]attr.Value{
							"error_code":      types.StringValue("IMAGE_NOT_FOUND"),
							"error_message":   types.StringValue("image missing"),
							"error_timestamp": types.StringValue("2024-01-02T03:04:05Z"),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenImageErrors(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
