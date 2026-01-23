// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack_theme

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

func TestExpandS3Location(t *testing.T) {
	t.Parallel()

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
				s3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringValue("my-bucket"),
					"s3_key":    types.StringValue("logo.png"),
				},
			),
			want: &awstypes.S3Location{
				S3Bucket: aws.String("my-bucket"),
				S3Key:    aws.String("logo.png"),
			},
		},
		{
			name: "bucket_only",
			input: types.ObjectValueMust(
				s3LocationObjectType.AttrTypes,
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
				s3LocationObjectType.AttrTypes,
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
			input:     types.ObjectUnknown(s3LocationObjectType.AttrTypes),
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
			t.Parallel()

			var diags diag.Diagnostics

			got := expandS3Location(ctx, tt.input, &diags)

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

func TestExpandFooterLinks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		set  types.Set
		want []awstypes.ThemeFooterLink
	}{
		{
			name: "null_set",
			set:  types.SetNull(footerLinkObjectType),
			want: nil,
		},
		{
			name: "single_footer_link",
			set: types.SetValueMust(
				footerLinkObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						footerLinkObjectType.AttrTypes,
						map[string]attr.Value{
							"display_name":    types.StringValue("Support"),
							"footer_link_url": types.StringValue("https://example.com/support"),
						},
					),
				},
			),
			want: []awstypes.ThemeFooterLink{
				{
					DisplayName:   aws.String("Support"),
					FooterLinkURL: aws.String("https://example.com/support"),
				},
			},
		},
		{
			name: "empty_set_returns_nil",
			set:  types.SetValueMust(footerLinkObjectType, []attr.Value{}),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := expandFooterLinks(ctx, tt.set, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
