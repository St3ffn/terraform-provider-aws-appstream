// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack_theme

import "github.com/hashicorp/terraform-plugin-framework/types"

type dataSourceModel struct {
	// ID is a synthetic identifier equal to the stack name.
	ID types.String `tfsdk:"id"`
	// StackName is the name of the AppStream stack (required).
	StackName types.String `tfsdk:"stack_name"`
	// TitleText is the browser tab title text (computed).
	TitleText types.String `tfsdk:"title_text"`
	// ThemeStyling is the color theme applied to the catalog page (computed).
	ThemeStyling types.String `tfsdk:"theme_styling"`
	// FooterLinks are optional footer links shown in the catalog page (computed).
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
