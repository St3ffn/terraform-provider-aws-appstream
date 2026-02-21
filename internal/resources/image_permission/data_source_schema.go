// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package image_permission

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
		Description: "Read AWS AppStream Image Permissions",
		MarkdownDescription: "Reads sharing permissions for a private AppStream image. " +
			"This data source can be used to inspect which AWS accounts a private image is shared with " +
			"and how those accounts are permitted to use the image.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "Name of the AppStream Image.",
				MarkdownDescription: "The name of the private AppStream image for which permissions are read.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"permissions": schema.SetNestedAttribute{
				Description:         "Image permissions per shared AWS account.",
				MarkdownDescription: "The list of AWS accounts the image is shared with and their corresponding permissions.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"shared_account_id": schema.StringAttribute{
							Description:         "Shared AWS account ID.",
							MarkdownDescription: "The AWS account ID with which the image is shared.",
							Computed:            true,
						},
						"image_permissions": schema.SingleNestedAttribute{
							Description:         "Permissions granted to the shared AWS account.",
							MarkdownDescription: "Defines how the shared AWS account is allowed to use the image.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"allow_fleet": schema.BoolAttribute{
									Description:         "Allow image usage for Fleets.",
									MarkdownDescription: "Whether the image can be used to create or update AppStream fleets.",
									Computed:            true,
								},
								"allow_image_builder": schema.BoolAttribute{
									Description:         "Allow image usage for Image Builders.",
									MarkdownDescription: "Whether the image can be used to create AppStream image builders.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}
