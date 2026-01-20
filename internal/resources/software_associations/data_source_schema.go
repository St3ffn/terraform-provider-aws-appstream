// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package software_associations

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (ds *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read AWS AppStream Software Associations",
		MarkdownDescription: "Reads license-included software associations for an AppStream image or image builder. " +
			"This data source is read-only and can be used to inspect which software is associated with a resource, " +
			"its deployment status, and any reported deployment errors.",
		Attributes: map[string]schema.Attribute{
			"associated_resource": schema.StringAttribute{
				Description:         "Associated AppStream resource ARN.",
				MarkdownDescription: "The ARN of the AppStream image or image builder for which software associations are read.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.Any(
						util.ValidARNWithServiceAndResource("appstream", "image/"),
						util.ValidARNWithServiceAndResource("appstream", "image-builder/"),
					),
				},
			},
			"software_associations": schema.SetNestedAttribute{
				Description:         "Software associations for the resource.",
				MarkdownDescription: "The list of license-included software applications associated with the resource and their deployment status.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"software_name": schema.StringAttribute{
							Description:         "Software name.",
							MarkdownDescription: "The name of the license-included application associated with the resource.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							Description:         "Software deployment status.",
							MarkdownDescription: "The deployment status of the license-included application, as reported by AWS.",
							Computed:            true,
						},
						"deployment_errors": schema.SetNestedAttribute{
							Description: "Deployment errors reported by AWS.",
							MarkdownDescription: "Informational list of error details for failed software deployments. " +
								"These errors are reported by AWS and do not affect Terraform lifecycle behavior.",
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"error_code": schema.StringAttribute{
										Description:         "Error code reported by AWS.",
										MarkdownDescription: "The error code reported by AWS for the software association.",
										Computed:            true,
									},
									"error_message": schema.StringAttribute{
										Description:         "Error message reported by AWS.",
										MarkdownDescription: "The human-readable error message reported by AWS for the software association.",
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
