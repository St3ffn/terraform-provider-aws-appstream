// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_theme

import "github.com/hashicorp/terraform-plugin-framework/types"

type resourceModel struct {
	// ID is a synthetic identifier equal to the stack name.
	ID types.String `tfsdk:"id"`
	// StackName is the name of the AppStream stack (required).
	StackName types.String `tfsdk:"stack_name"`
	// TitleText is the browser tab title text (required).
	TitleText types.String `tfsdk:"title_text"`
	// ThemeStyling is the color theme applied to the catalog page (required).
	ThemeStyling types.String `tfsdk:"theme_styling"`
	// OrganizationLogoS3Location is the S3 location of the organization logo (required).
	OrganizationLogoS3Location types.Object `tfsdk:"organization_logo_s3_location"`
	// FaviconS3Location is the S3 location of the favicon (required).
	FaviconS3Location types.Object `tfsdk:"favicon_s3_location"`
	// FooterLinks are optional footer links shown in the catalog page (optional).
	FooterLinks types.Set `tfsdk:"footer_links"`
	// State specifies whether the theme is enabled or disabled (computed).
	State types.String `tfsdk:"state"`
	// CreatedTime is the timestamp when the theme was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
	// ThemeOrganizationLogoURL is the URL of the organization logo displayed in the catalog page header (computed).
	ThemeOrganizationLogoURL types.String `tfsdk:"theme_organization_logo_url"`
	// ThemeFaviconURL is the URL of the favicon displayed in the browser tab during streaming sessions (computed).
	ThemeFaviconURL types.String `tfsdk:"theme_favicon_url"`
}

type s3LocationModel struct {
	// S3Bucket is the name of the S3 bucket.
	S3Bucket types.String `tfsdk:"s3_bucket"`
	// S3Key is the S3 object key.
	S3Key types.String `tfsdk:"s3_key"`
}

type footerLinkModel struct {
	// DisplayName is the name displayed for the footer link (required).
	DisplayName types.String `tfsdk:"display_name"`
	// FooterLinkURL is the URL that the footer link points to (required).
	FooterLinkURL types.String `tfsdk:"footer_link_url"`
}
