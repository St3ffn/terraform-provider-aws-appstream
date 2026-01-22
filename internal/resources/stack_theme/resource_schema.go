// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package stack_theme

import (
	"context"
	"regexp"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream Stack Theme",
		MarkdownDescription: "Manages a custom branding theme for an AppStream stack. " +
			"A stack theme customizes the appearance of the streaming application catalog page, " +
			"including colors, title text, footer links, and branding assets such as logos and favicons.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream stack theme.",
				MarkdownDescription: "A synthetic identifier for the stack theme, equal to the stack name. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"stack_name": schema.StringAttribute{
				Description: "Name of the AppStream stack.",
				MarkdownDescription: "The name of the AppStream stack that the theme is applied to. " +
					"Changing this value forces the theme to be replaced.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"title_text": schema.StringAttribute{
				Description:         "Browser tab title text.",
				MarkdownDescription: "The title displayed at the top of the browser tab during application streaming sessions.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 300),
				},
			},
			"theme_styling": schema.StringAttribute{
				Description:         "Theme color styling.",
				MarkdownDescription: "The color theme applied to links, text, buttons, and accents in the streaming catalog page.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						util.AWSEnumToSlice(awstypes.ThemeStyling.Values)...,
					),
				},
			},
			"organization_logo_s3_location": schema.SingleNestedAttribute{
				Description: "Organization logo S3 location.",
				MarkdownDescription: "Specifies the Amazon S3 location of the organization logo displayed " +
					"in the header of the streaming application catalog page.",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"s3_bucket": schema.StringAttribute{
						Description:         "S3 bucket name.",
						MarkdownDescription: "The name of the Amazon S3 bucket.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(3, 63),
						},
					},
					"s3_key": schema.StringAttribute{
						Description:         "S3 object key.",
						MarkdownDescription: "The S3 object key of the organization logo.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 1024),
						},
					},
				},
			},
			"favicon_s3_location": schema.SingleNestedAttribute{
				Description: "Favicon S3 location.",
				MarkdownDescription: "Specifies the Amazon S3 location of the favicon displayed in browser tabs " +
					"during application streaming sessions.",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"s3_bucket": schema.StringAttribute{
						Description:         "S3 bucket name.",
						MarkdownDescription: "The name of the Amazon S3 bucket.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(3, 63),
						},
					},
					"s3_key": schema.StringAttribute{
						Description:         "S3 object key.",
						MarkdownDescription: "The S3 object key of the favicon.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 1024),
						},
					},
				},
			},
			"footer_links": schema.SetNestedAttribute{
				Description: "Footer links for the streaming catalog page.",
				MarkdownDescription: "Optional links displayed in the footer of the streaming application catalog page, " +
					"such as support, documentation, or marketing links.",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.StringAttribute{
							Description:         "Display name of the footer link.",
							MarkdownDescription: "The name displayed for the footer link in the catalog page.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 300),
							},
						},
						"footer_link_url": schema.StringAttribute{
							Description:         "URL of the footer link.",
							MarkdownDescription: "The URL that the footer link points to.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 1000),
							},
						},
					},
				},
			},
			"state": schema.StringAttribute{
				Description:         "State of the stack theme.",
				MarkdownDescription: "The current state of the AppStream stack theme.",
				Computed:            true,
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the stack theme was created.",
				MarkdownDescription: "The timestamp when the stack theme was created, in RFC 3339 format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"theme_organization_logo_url": schema.StringAttribute{
				Description:         "Organization logo URL.",
				MarkdownDescription: "The URL of the logo that is displayed in the header of the streaming application catalog page.",
				Computed:            true,
			},
			"theme_favicon_url": schema.StringAttribute{
				Description:         "Favicon URL.",
				MarkdownDescription: "The URL of the icon that is displayed at the top of the browser tab during application streaming sessions.",
				Computed:            true,
			},
		},
	}
}
