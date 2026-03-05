// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Schema returns the static Terraform schema for this component.
// It defines supported attributes/blocks plus validators and plan behavior.
func (ds *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read an AWS AppStream App Block Builder",
		MarkdownDescription: "Reads an AppStream app block builder. " +
			"This data source can be used to reference an existing AppStream app block builder that is managed outside of Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the AppStream App Block Builder.",
				MarkdownDescription: "A synthetic identifier for the app block builder, equal to the app block builder name.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the AppStream App Block Builder.",
				MarkdownDescription: "The name of the AppStream app block builder to read.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"instance_type": schema.StringAttribute{
				Description:         "Instance type for the App Block Builder.",
				MarkdownDescription: "The instance type used when launching the app block builder.",
				Computed:            true,
			},
			"platform": schema.StringAttribute{
				Description:         "App block builder platform.",
				MarkdownDescription: "The operating system platform of the app block builder.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the AppStream App Block Builder.",
				MarkdownDescription: "The app block builder description, if set.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				Description:         "Display name of the AppStream App Block Builder.",
				MarkdownDescription: "The display name of the app block builder shown in the AppStream user interface.",
				Computed:            true,
			},
			"disable_imds_v1": schema.BoolAttribute{
				Description: "Disable IMDSv1 and enforce IMDSv2 for the AppStream App Block Builder.",
				MarkdownDescription: "Whether Instance Metadata Service Version 1 (IMDSv1) is disabled for the " +
					"AppStream app block builder. If `true`, only IMDSv2 is enabled. " +
					"If `false`, both IMDSv1 and IMDSv2 are enabled.",
				Computed: true,
			},
			"vpc_config": schema.SingleNestedAttribute{
				Description:         "VPC configuration for the App Block Builder.",
				MarkdownDescription: "The VPC configuration used by the app block builder.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"subnet_ids": schema.SetAttribute{
						Description:         "Subnet IDs.",
						MarkdownDescription: "The subnet IDs in which the app block builder is launched.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"security_group_ids": schema.SetAttribute{
						Description:         "Security group IDs.",
						MarkdownDescription: "The security group IDs associated with the app block builder.",
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"iam_role_arn": schema.StringAttribute{
				Description:         "IAM role ARN.",
				MarkdownDescription: "The ARN of the IAM role applied to the app block builder.",
				Computed:            true,
			},
			"enable_default_internet_access": schema.BoolAttribute{
				Description:         "Enable default internet access.",
				MarkdownDescription: "Whether the app block builder has access to the internet.",
				Computed:            true,
			},
			"access_endpoints": schema.SetNestedAttribute{
				Description:         "VPC access endpoints for the App Block Builder.",
				MarkdownDescription: "Interface VPC endpoints through which administrators can connect to the app block builder.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"endpoint_type": schema.StringAttribute{
							Description:         "Endpoint type.",
							MarkdownDescription: "The type of interface endpoint.",
							Computed:            true,
						},
						"vpce_id": schema.StringAttribute{
							Description:         "VPC endpoint ID.",
							MarkdownDescription: "The identifier of the interface VPC endpoint.",
							Computed:            true,
						},
					},
				},
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the App Block Builder.",
				MarkdownDescription: "Tags assigned to the app block builder.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream App Block Builder.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream app block builder.",
				Computed:            true,
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the App Block Builder was created.",
				MarkdownDescription: "The timestamp when the app block builder was created, in RFC 3339 format.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				Description:         "State of the AppStream App Block Builder.",
				MarkdownDescription: "The state of the AppStream app block builder.",
				Computed:            true,
			},
			"state_change_reason": schema.SingleNestedAttribute{
				Description:         "State change reason.",
				MarkdownDescription: "The reason for the most recent app block builder state change, if applicable.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"code": schema.StringAttribute{
						Description:         "State change reason code.",
						MarkdownDescription: "The code describing why the app block builder state changed.",
						Computed:            true,
					},
					"message": schema.StringAttribute{
						Description:         "State change reason message.",
						MarkdownDescription: "The human-readable message describing the state change.",
						Computed:            true,
					},
				},
			},
			"app_block_builder_errors": schema.SetNestedAttribute{
				Description:         "Errors reported by AWS for the App Block Builder.",
				MarkdownDescription: "Informational list of errors reported by AWS for the app block builder.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							Description:         "Error code reported by AWS.",
							MarkdownDescription: "The error code reported by AWS.",
							Computed:            true,
						},
						"error_message": schema.StringAttribute{
							Description:         "Error message reported by AWS.",
							MarkdownDescription: "The human-readable error message reported by AWS.",
							Computed:            true,
						},
						"error_timestamp": schema.StringAttribute{
							Description:         "Error timestamp.",
							MarkdownDescription: "The timestamp when the error occurred, in RFC 3339 format.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
