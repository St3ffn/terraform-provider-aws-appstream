// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream App Block Builder",
		MarkdownDescription: "Manages an AppStream app block builder. " +
			"An app block builder is a reusable resource used to package and test AppStream app blocks by installing " +
			"and configuring applications before associating them with Elastic fleets.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream app block builder.",
				MarkdownDescription: "A synthetic identifier for the app block builder, equal to the app block builder name. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the AppStream app block builder.",
				MarkdownDescription: "The name of the AppStream app block builder. " +
					"Changing this value forces the app block builder to be replaced.",
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
			"instance_type": schema.StringAttribute{
				Description:         "Instance type for the app block builder.",
				MarkdownDescription: "The instance type used when launching the app block builder.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"platform": schema.StringAttribute{
				Description:         "App block builder platform.",
				MarkdownDescription: "The operating system platform of the app block builder.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description of the AppStream App Block Builder.",
				MarkdownDescription: "The app block builder description, if set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"display_name": schema.StringAttribute{
				Description:         "Display name of the app block builder.",
				MarkdownDescription: "The display name of the app block builder shown in the AppStream user interface.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"vpc_config": schema.SingleNestedAttribute{
				Description: "VPC configuration for the app block builder.",
				MarkdownDescription: "The VPC configuration used by the app block builder. " +
					"App block builders require at least two subnet IDs in different Availability Zones.",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"subnet_ids": schema.SetAttribute{
						Description: "Subnet IDs.",
						MarkdownDescription: "The subnet IDs in which the app block builder is launched. " +
							"Specify at least two subnet IDs in different availability zones.",
						Required:    true,
						ElementType: types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
						},
					},
					"security_group_ids": schema.SetAttribute{
						Description:         "Security group IDs.",
						MarkdownDescription: "The security group IDs associated with the app block builder.",
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
				MarkdownDescription: "The ARN of the IAM role applied to the app block builder.",
				Optional:            true,
				Validators: []validator.String{
					util.ValidARNWithServiceAndResource("iam", "role/"),
				},
			},
			"enable_default_internet_access": schema.BoolAttribute{
				Description:         "Enable default internet access.",
				MarkdownDescription: "Whether the app block builder has access to the internet.",
				Optional:            true,
				Computed:            true,
			},
			"access_endpoints": schema.SetNestedAttribute{
				Description:         "VPC access endpoints for the app block builder.",
				MarkdownDescription: "Interface VPC endpoints through which administrators can connect to the app block builder.",
				Optional:            true,
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
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the AppStream app block builder.",
				MarkdownDescription: "A map of tags assigned to the AppStream app block builder.",
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
				Description:         "ARN of the AppStream app block builder.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream app block builder.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the app block builder was created.",
				MarkdownDescription: "The timestamp when the app block builder was created, in RFC 3339 format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description:         "State of the AppStream app block builder.",
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
				Description: "Errors reported by AWS for the app block builder.",
				MarkdownDescription: "Informational list of errors reported by AWS for the app block builder. " +
					"These errors do not affect Terraform lifecycle behavior.",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							Description:         "Error code reported by AWS.",
							MarkdownDescription: "The error code reported by AWS for the app block builder.",
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
