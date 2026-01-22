// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack_theme

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (ds *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read an AWS AppStream Stack Theme",
		MarkdownDescription: "Reads a custom branding theme for an AppStream stack. " +
			"This data source can be used to reference an existing AppStream stack theme " +
			"that is managed outside of Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the AppStream stack theme.",
				MarkdownDescription: "A synthetic identifier for the stack theme, equal to the stack name.",
				Computed:            true,
			},
			"stack_name": schema.StringAttribute{
				Description:         "Name of the AppStream stack.",
				MarkdownDescription: "The name of the AppStream stack whose theme is read.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"title_text": schema.StringAttribute{
				Description:         "Browser tab title text.",
				MarkdownDescription: "The title displayed at the top of the browser tab during streaming sessions.",
				Computed:            true,
			},
			"theme_styling": schema.StringAttribute{
				Description:         "Theme color styling.",
				MarkdownDescription: "The color theme applied to links, text, buttons, and accents in the catalog page.",
				Computed:            true,
			},
			"footer_links": schema.SetNestedAttribute{
				Description:         "Footer links for the catalog page.",
				MarkdownDescription: "Links displayed in the footer of the streaming application catalog page.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.StringAttribute{
							Description:         "Display name of the footer link.",
							MarkdownDescription: "The name displayed for the footer link.",
							Computed:            true,
						},
						"footer_link_url": schema.StringAttribute{
							Description:         "URL of the footer link.",
							MarkdownDescription: "The URL that the footer link points to.",
							Computed:            true,
						},
					},
				},
			},
			"state": schema.StringAttribute{
				Description:         "State of the stack theme.",
				MarkdownDescription: "The current state of the AppStream stack theme.",
				Computed:            true,
			},
			"theme_organization_logo_url": schema.StringAttribute{
				Description:         "Organization logo URL.",
				MarkdownDescription: "The URL of the logo displayed in the catalog page header.",
				Computed:            true,
			},
			"theme_favicon_url": schema.StringAttribute{
				Description:         "Favicon URL.",
				MarkdownDescription: "The URL of the icon displayed in the browser tab during streaming sessions.",
				Computed:            true,
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the stack theme was created.",
				MarkdownDescription: "The timestamp when the stack theme was created, in RFC 3339 format.",
				Computed:            true,
			},
		},
	}
}
