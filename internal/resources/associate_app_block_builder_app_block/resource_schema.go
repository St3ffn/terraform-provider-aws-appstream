// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_app_block_builder_app_block

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Schema returns the static Terraform schema for this component.
// It defines supported attributes/blocks plus validators and plan behavior.
func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream App Block Builder–App Block Association",
		MarkdownDescription: "Manages the association between an AppStream app block builder and an AppStream app block. " +
			"This resource represents the relationship only and does not create or manage the underlying app block builder or app block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream App Block Builder–App Block association.",
				MarkdownDescription: "A synthetic identifier for the association, composed of the app block builder name and app block ARN. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_block_builder_name": schema.StringAttribute{
				Description: "Name of the AppStream App Block Builder.",
				MarkdownDescription: "The name of the AppStream app block builder to associate with the app block. " +
					"Changing this value forces the association to be replaced.",
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
			"app_block_arn": schema.StringAttribute{
				Description: "ARN of the AppStream App Block.",
				MarkdownDescription: "The ARN of the AppStream app block to associate with the app block builder. " +
					"Changing this value forces the association to be replaced.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					util.ValidARNWithServiceAndResource("appstream", "app-block/"),
				},
			},
		},
	}
}
