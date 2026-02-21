// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package image_permission

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage AWS AppStream Image Sharing Permission",
		MarkdownDescription: "Manages sharing permissions for a private AppStream image. " +
			"Image permissions control whether a private image can be used by fleets or image builders " +
			"in another AWS account. " +
			"Each permission grants access to **one AWS account** and defines whether the image is usable " +
			"for fleets, image builders, or both.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream Image Permission.",
				MarkdownDescription: "A synthetic identifier for the image permission, composed of the image name " +
					"and the shared AWS account ID. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the AppStream Image.",
				MarkdownDescription: "The name of the private AppStream image for which permissions are managed. " +
					"Changing this value forces the image permission to be replaced.",
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
			"shared_account_id": schema.StringAttribute{
				Description: "AWS account ID to share the image with.",
				MarkdownDescription: "The 12-digit AWS account ID that is granted permissions to use the image. " +
					"Each image permission resource manages access for exactly one AWS account.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d+$`),
						"must match ^\\d+$",
					),
				},
			},
			"image_permissions": schema.SingleNestedAttribute{
				Description: "Permissions granted to the shared AWS account.",
				MarkdownDescription: "Defines how the shared AWS account is allowed to use the image. " +
					"At least one permission must be enabled.",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"allow_fleet": schema.BoolAttribute{
						Description: "Allow image usage for Fleets.",
						MarkdownDescription: "Whether the image can be used to create or update AppStream fleets " +
							"in the shared AWS account.",
						Required: true,
					},
					"allow_image_builder": schema.BoolAttribute{
						Description: "Allow image usage for Image Builders.",
						MarkdownDescription: "Whether the image can be used to create AppStream image builders " +
							"in the shared AWS account.",
						Required: true,
					},
				},
			},
		},
	}
}
