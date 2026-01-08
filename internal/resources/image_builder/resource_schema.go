// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream Image Builder",
		MarkdownDescription: "Manages an AppStream image builder. " +
			"An image builder is used to create custom AppStream images by installing applications " +
			"and capturing the resulting configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream image builder.",
				MarkdownDescription: "A synthetic identifier for the image builder, equal to the image builder name. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the AppStream image builder.",
				MarkdownDescription: "The name of the AppStream image builder. " +
					"Changing this value forces the image builder to be replaced.",
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
			"image_name": schema.StringAttribute{
				Description: "Name of the AppStream image.",
				MarkdownDescription: "The name of the AppStream image used to create the image builder. " +
					"Either `image_name` or `image_arn` must be specified.",
				Optional: true,
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
			"image_arn": schema.StringAttribute{
				Description: "ARN of the AppStream image.",
				MarkdownDescription: "The ARN of the AppStream image used to create the image builder. " +
					"Either `image_name` or `image_arn` must be specified.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					util.ValidARNWithServiceAndResource("appstream", "image/"),
				},
			},
			"instance_type": schema.StringAttribute{
				Description:         "Instance type for the image builder.",
				MarkdownDescription: "The instance type used when launching the image builder.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description of the AppStream Image Builder.",
				MarkdownDescription: "The image builder description, if set.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"display_name": schema.StringAttribute{
				Description:         "Display name of the image builder.",
				MarkdownDescription: "The display name of the image builder shown in the AppStream user interface.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"vpc_config": schema.SingleNestedAttribute{
				Description: "VPC configuration for the image builder.",
				MarkdownDescription: "The VPC configuration used by the image builder. " +
					"Image builders use exactly one subnet.",
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"subnet_ids": schema.SetAttribute{
						Description:         "Subnet IDs.",
						MarkdownDescription: "The subnet IDs in which the image builder is launched.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"security_group_ids": schema.SetAttribute{
						Description:         "Security group IDs.",
						MarkdownDescription: "The security group IDs associated with the image builder.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtMost(5),
							setvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(1),
							),
						},
					},
				},
			},
			"iam_role_arn": schema.StringAttribute{
				Description:         "IAM role ARN.",
				MarkdownDescription: "The ARN of the IAM role applied to the image builder.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					util.ValidARNWithServiceAndResource("iam", "role/"),
				},
			},
			"enable_default_internet_access": schema.BoolAttribute{
				Description:         "Enable default internet access.",
				MarkdownDescription: "Whether the image builder has access to the internet.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"domain_join_info": schema.SingleNestedAttribute{
				Description: "Active Directory domain join configuration.",
				MarkdownDescription: "Specifies the Active Directory domain and organizational unit used to join " +
					"the image builder to a Microsoft Active Directory domain.",
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"directory_name": schema.StringAttribute{
						Description:         "Directory name.",
						MarkdownDescription: "The fully qualified domain name of the Active Directory.",
						Optional:            true,
					},
					"organizational_unit_distinguished_name": schema.StringAttribute{
						Description:         "Organizational unit distinguished name.",
						MarkdownDescription: "The distinguished name of the organizational unit for computer accounts.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(2000),
						},
					},
				},
			},
			"appstream_agent_version": schema.StringAttribute{
				Description:         "AppStream agent version.",
				MarkdownDescription: "The AppStream agent version used by the image builder.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"access_endpoints": schema.SetNestedAttribute{
				Description:         "VPC access endpoints for the image builder.",
				MarkdownDescription: "Interface VPC endpoints through which administrators can connect to the image builder.",
				Optional:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 4),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"endpoint_type": schema.StringAttribute{
							Description:         "Endpoint type.",
							MarkdownDescription: "The type of interface endpoint.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("STREAMING"),
							},
						},
						"vpce_id": schema.StringAttribute{
							Description:         "VPC endpoint ID.",
							MarkdownDescription: "The identifier of the interface VPC endpoint.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
			"root_volume_config": schema.SingleNestedAttribute{
				Description:         "Root volume configuration.",
				MarkdownDescription: "Specifies the root volume configuration for the image builder.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"volume_size_in_gb": schema.Int32Attribute{
						Description:         "Root volume size.",
						MarkdownDescription: "The size of the root volume, in GiB.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.Int32{
							int32validator.Between(200, 500),
						},
					},
				},
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the AppStream image builder.",
				MarkdownDescription: "A map of tags assigned to the AppStream image builder.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.Map{
					mapvalidator.SizeAtMost(50),
					mapvalidator.KeysAre(
						stringvalidator.LengthBetween(1, 128),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[\p{L}\p{Z}\p{N}_.:/=+\-@]*$`),
							"must match ^[\\p{L}\\p{Z}\\p{N}_.:/=+\\-@]*$",
						),
					),
					mapvalidator.ValueStringsAre(
						stringvalidator.LengthAtMost(256),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^([\p{L}\p{Z}\p{N}_.:/=+\-@]*)$`),
							"must match ^([\\p{L}\\p{Z}\\p{N}_.:/=+\\-@]*)$",
						),
					),
				},
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream image builder.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream image builder.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the image builder was created.",
				MarkdownDescription: "The timestamp when the image builder was created, in RFC 3339 format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Description: "Errors reported by AWS for the image builder.",
				MarkdownDescription: "Informational list of errors reported by AWS for the image builder. " +
					"These errors do not affect Terraform lifecycle behavior.",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							Description:         "Error code reported by AWS.",
							MarkdownDescription: "The error code reported by AWS for the image builder.",
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
