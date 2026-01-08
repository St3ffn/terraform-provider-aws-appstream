// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (ds *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read an AWS AppStream Image Builder",
		MarkdownDescription: "Reads an AppStream image builder. " +
			"This data source can be used to reference an existing AppStream image builder that is managed outside of Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the AppStream image builder.",
				MarkdownDescription: "A synthetic identifier for the image builder, equal to the image builder name.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the AppStream image builder.",
				MarkdownDescription: "The name of the AppStream image builder to read.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"image_arn": schema.StringAttribute{
				Description:         "ARN of the AppStream image.",
				MarkdownDescription: "The ARN of the AppStream image used to create the image builder.",
				Computed:            true,
			},
			"instance_type": schema.StringAttribute{
				Description:         "Instance type for the image builder.",
				MarkdownDescription: "The instance type used when launching the image builder.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the AppStream image builder.",
				MarkdownDescription: "The image builder description, if set.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				Description:         "Display name of the AppStream image builder.",
				MarkdownDescription: "The name displayed to users in the AppStream user interface.",
				Computed:            true,
			},
			"vpc_config": schema.SingleNestedAttribute{
				Description:         "VPC configuration for the image builder.",
				MarkdownDescription: "The VPC configuration used by the image builder, if applicable.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"subnet_ids": schema.SetAttribute{
						Description:         "Subnet IDs.",
						MarkdownDescription: "The subnet IDs in which the image builder is launched.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"security_group_ids": schema.SetAttribute{
						Description:         "Security group IDs.",
						MarkdownDescription: "The security group IDs associated with the image builder.",
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"iam_role_arn": schema.StringAttribute{
				Description:         "IAM role ARN.",
				MarkdownDescription: "The ARN of the IAM role applied to the image builder.",
				Computed:            true,
			},
			"enable_default_internet_access": schema.BoolAttribute{
				Description:         "Enable default internet access.",
				MarkdownDescription: "Whether the image builder has access to the internet.",
				Computed:            true,
			},
			"domain_join_info": schema.SingleNestedAttribute{
				Description:         "Active Directory domain join configuration.",
				MarkdownDescription: "The Active Directory configuration used by the image builder, if applicable.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"directory_name": schema.StringAttribute{
						Description:         "Directory name.",
						MarkdownDescription: "The fully qualified domain name of the Active Directory.",
						Computed:            true,
					},
					"organizational_unit_distinguished_name": schema.StringAttribute{
						Description:         "Organizational unit distinguished name.",
						MarkdownDescription: "The distinguished name of the organizational unit for computer accounts.",
						Computed:            true,
					},
				},
			},
			"appstream_agent_version": schema.StringAttribute{
				Description:         "AppStream agent version.",
				MarkdownDescription: "The AppStream agent version used by the image builder.",
				Computed:            true,
			},
			"access_endpoints": schema.SetNestedAttribute{
				Description:         "VPC access endpoints for the image builder.",
				MarkdownDescription: "Interface VPC endpoints through which administrators can connect to the image builder.",
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
			"root_volume_config": schema.SingleNestedAttribute{
				Description:         "Root volume configuration.",
				MarkdownDescription: "The root volume configuration of the image builder, if set.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"volume_size_in_gb": schema.Int32Attribute{
						Description:         "Root volume size.",
						MarkdownDescription: "The size of the root volume, in GiB.",
						Computed:            true,
					},
				},
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the AppStream image builder.",
				MarkdownDescription: "Tags assigned to the AppStream image builder.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream image builder.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream image builder.",
				Computed:            true,
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the image builder was created.",
				MarkdownDescription: "The timestamp when the image builder was created, in RFC 3339 format.",
				Computed:            true,
			},
			"platform": schema.StringAttribute{
				Description:         "Image builder platform.",
				MarkdownDescription: "The operating system platform of the image builder.",
				Computed:            true,
			},
			"network_access_configuration": schema.SingleNestedAttribute{
				Description:         "Network access configuration.",
				MarkdownDescription: "Network details of the elastic network interface attached to the image builder.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"eni_private_ip_address": schema.StringAttribute{
						Description:         "Private IP address.",
						MarkdownDescription: "The private IP address of the elastic network interface.",
						Computed:            true,
					},
					"eni_ipv6_addresses": schema.SetAttribute{
						Description:         "IPv6 addresses.",
						MarkdownDescription: "The IPv6 addresses assigned to the elastic network interface.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"eni_id": schema.StringAttribute{
						Description:         "Elastic network interface ID.",
						MarkdownDescription: "The identifier of the elastic network interface.",
						Computed:            true,
					},
				},
			},
			"latest_appstream_agent_version": schema.StringAttribute{
				Description:         "Latest AppStream agent version indicator.",
				MarkdownDescription: "Indicates whether the image builder is using the latest AppStream agent version.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				Description:         "State of the AppStream image builder.",
				MarkdownDescription: "The state of the AppStream image builder.",
				Computed:            true,
			},
			"state_change_reason": schema.SingleNestedAttribute{
				Description:         "State change reason.",
				MarkdownDescription: "The reason for the most recent image builder state change, if applicable.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"code": schema.StringAttribute{
						Description:         "State change reason code.",
						MarkdownDescription: "The code describing why the image builder state changed.",
						Computed:            true,
					},
					"message": schema.StringAttribute{
						Description:         "State change reason message.",
						MarkdownDescription: "The human-readable message describing the state change.",
						Computed:            true,
					},
				},
			},
			"image_builder_errors": schema.SetNestedAttribute{
				Description:         "Errors reported by AWS for the image builder.",
				MarkdownDescription: "Informational list of errors reported by AWS for the image builder.",
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
