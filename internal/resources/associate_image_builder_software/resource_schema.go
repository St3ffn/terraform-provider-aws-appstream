// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Schema returns the static Terraform schema for this component.
// It defines supported attributes/blocks plus validators and plan behavior.
func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream Image Builder Software Association",
		MarkdownDescription: "Manages the association between an AppStream image builder and one or more " +
			"license-included software packages. " +
			"This resource represents the relationship only and does not create or manage the underlying " +
			"image builder or software packages. " +
			"Optionally, a software deployment can be triggered after the association is created. " +
			"After `terraform import`, set `software_names` (and optionally `deploy`) in configuration because " +
			"AWS read APIs cannot reconstruct Terraform ownership intent for those attributes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the Image Builder software association.",
				MarkdownDescription: "A synthetic identifier for the association, equal to the image builder ARN. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"image_builder_arn": schema.StringAttribute{
				Description: "ARN of the AppStream Image Builder.",
				MarkdownDescription: "The ARN of the AppStream image builder to which the software is associated. " +
					"Changing this value forces the association to be replaced.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					util.ValidARNWithServiceAndResource("appstream", "image-builder/"),
				},
			},
			"software_names": schema.SetAttribute{
				Description: "Software names associated with the Image Builder.",
				MarkdownDescription: "A set of license-included software package names to associate with the AppStream " +
					"image builder. Changes to this set result in software being associated or disassociated accordingly. " +
					"Required for normal management and post-import reconciliation because AWS does not return the " +
					"Terraform-managed target set.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"deploy": schema.BoolAttribute{
				Description: "Trigger software deployment.",
				MarkdownDescription: "Whether to trigger a software deployment to the image builder after associating " +
					"the software. When set to `true`, a deployment is started during apply. " +
					"This setting does not track deployment progress or completion. " +
					"This attribute is Terraform intent and is not derivable from AWS during import.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"associations": schema.MapNestedAttribute{
				Description: "Software association details.",
				MarkdownDescription: "Per-software association status and deployment details as reported by AWS. " +
					"This attribute is informational only and does not affect Terraform lifecycle behavior. " +
					"On import, values are filtered by configured `software_names`.",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"status": schema.StringAttribute{
							Description: "Software deployment status.",
							MarkdownDescription: "The deployment status of the license-included application for the image builder, " +
								"as reported by AWS.",
							Computed: true,
						},
						"deployment_errors": schema.SetNestedAttribute{
							Description: "Deployment errors reported by AWS.",
							MarkdownDescription: "Informational list of error details for failed software deployments. " +
								"These errors do not affect Terraform lifecycle behavior.",
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"error_code": schema.StringAttribute{
										Description:         "Error code reported by AWS.",
										MarkdownDescription: "The error code reported by AWS for the association.",
										Computed:            true,
									},
									"error_message": schema.StringAttribute{
										Description:         "Error message reported by AWS.",
										MarkdownDescription: "The human-readable error message reported by AWS for the association.",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
