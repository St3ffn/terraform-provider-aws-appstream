// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_theme

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenS3Location(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		in   *awstypes.S3Location
		want types.Object
	}{
		{
			name: "nil_input",
			in:   nil,
			want: types.ObjectNull(s3LocationObjectType.AttrTypes),
		},
		{
			name: "bucket_and_key_set",
			in: &awstypes.S3Location{
				S3Bucket: aws.String("my-bucket"),
				S3Key:    aws.String("logo.png"),
			},
			want: types.ObjectValueMust(
				s3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringValue("my-bucket"),
					"s3_key":    types.StringValue("logo.png"),
				},
			),
		},
		{
			name: "partial_values",
			in: &awstypes.S3Location{
				S3Bucket: aws.String("my-bucket"),
				S3Key:    nil,
			},
			want: types.ObjectValueMust(
				s3LocationObjectType.AttrTypes,
				map[string]attr.Value{
					"s3_bucket": types.StringValue("my-bucket"),
					"s3_key":    types.StringNull(),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := flattenS3Location(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlattenFooterLinks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		in   []awstypes.ThemeFooterLink
		want types.Set
	}{
		{
			name: "nil_slice_returns_null",
			in:   nil,
			want: types.SetNull(footerLinkObjectType),
		},
		{
			name: "empty_slice_returns_null",
			in:   []awstypes.ThemeFooterLink{},
			want: types.SetNull(footerLinkObjectType),
		},
		{
			name: "single_footer_link",
			in: []awstypes.ThemeFooterLink{
				{
					DisplayName:   aws.String("Support"),
					FooterLinkURL: aws.String("https://example.com/support"),
				},
			},
			want: types.SetValueMust(
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
		},
		{
			name: "partial_values",
			in: []awstypes.ThemeFooterLink{
				{
					DisplayName:   nil,
					FooterLinkURL: aws.String("https://example.com"),
				},
			},
			want: types.SetValueMust(
				footerLinkObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						footerLinkObjectType.AttrTypes,
						map[string]attr.Value{
							"display_name":    types.StringNull(),
							"footer_link_url": types.StringValue("https://example.com"),
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

			got := flattenFooterLinks(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
