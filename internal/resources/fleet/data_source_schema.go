// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package fleet

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
		Description: "Read an AWS AppStream Fleet",
		MarkdownDescription: "Reads an AppStream fleet. " +
			"This data source can be used to reference an existing AppStream fleet that is managed outside of Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the AppStream fleet.",
				MarkdownDescription: "A synthetic identifier for the fleet, equal to the fleet name.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the AppStream fleet.",
				MarkdownDescription: "The name of the AppStream fleet to read.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"image_name": schema.StringAttribute{
				Description:         "Name of the AppStream image.",
				MarkdownDescription: "The name of the AppStream image used to create the fleet, if set.",
				Computed:            true,
			},
			"image_arn": schema.StringAttribute{
				Description:         "ARN of the AppStream image.",
				MarkdownDescription: "The ARN of the AppStream image used to create the fleet, if set.",
				Computed:            true,
			},
			"instance_type": schema.StringAttribute{
				Description:         "EC2 instance type for fleet instances.",
				MarkdownDescription: "The EC2 instance type used by the fleet.",
				Computed:            true,
			},
			"fleet_type": schema.StringAttribute{
				Description:         "Type of the AppStream fleet.",
				MarkdownDescription: "The fleet type. Valid values are `ON_DEMAND`, `ALWAYS_ON`, or `ELASTIC`.",
				Computed:            true,
			},
			"compute_capacity": schema.SingleNestedAttribute{
				Description: "Compute capacity configuration.",
				MarkdownDescription: "Describes the compute capacity configuration of the fleet, if applicable. " +
					"For non-elastic fleets, either instance-based or session-based capacity is returned.",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"desired_instances": schema.Int32Attribute{
						Description:         "Desired number of fleet instances.",
						MarkdownDescription: "The desired number of streaming instances for a single-session fleet.",
						Computed:            true,
					},
					"desired_sessions": schema.Int32Attribute{
						Description:         "Desired number of streaming sessions.",
						MarkdownDescription: "The desired number of concurrent streaming sessions for a multi-session fleet.",
						Computed:            true,
					},
				},
			},
			"vpc_config": schema.SingleNestedAttribute{
				Description: "VPC configuration for the fleet.",
				MarkdownDescription: "The VPC configuration used by the fleet, if applicable. " +
					"This configuration is present for elastic fleets.",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"subnet_ids": schema.SetAttribute{
						Description:         "Subnet IDs.",
						MarkdownDescription: "The subnet IDs associated with the fleet.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"security_group_ids": schema.SetAttribute{
						Description:         "Security group IDs.",
						MarkdownDescription: "The security group IDs associated with the fleet.",
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"max_user_duration_in_seconds": schema.Int32Attribute{
				Description:         "Maximum user session duration.",
				MarkdownDescription: "The maximum length of time that a streaming session can remain active.",
				Computed:            true,
			},
			"disconnect_timeout_in_seconds": schema.Int32Attribute{
				Description:         "Session disconnect timeout.",
				MarkdownDescription: "The amount of time that a disconnected session is allowed to remain active.",
				Computed:            true,
			},
			"idle_disconnect_timeout_in_seconds": schema.Int32Attribute{
				Description:         "Idle session disconnect timeout.",
				MarkdownDescription: "The amount of time, in seconds, that a session can remain idle before being disconnected.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the AppStream fleet.",
				MarkdownDescription: "The fleet description, if set.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				Description:         "Display name of the AppStream fleet.",
				MarkdownDescription: "The name displayed to users in the AppStream user interface.",
				Computed:            true,
			},
			"enable_default_internet_access": schema.BoolAttribute{
				Description:         "Enable default internet access.",
				MarkdownDescription: "Whether instances in the fleet have access to the internet.",
				Computed:            true,
			},
			"domain_join_info": schema.SingleNestedAttribute{
				Description:         "Active Directory domain join configuration.",
				MarkdownDescription: "The Active Directory configuration used by the fleet, if applicable.",
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
			"iam_role_arn": schema.StringAttribute{
				Description:         "IAM role ARN.",
				MarkdownDescription: "The ARN of the IAM role applied to fleet instances.",
				Computed:            true,
			},
			"stream_view": schema.StringAttribute{
				Description:         "Streaming view configuration.",
				MarkdownDescription: "Controls which streaming protocol views are enabled.",
				Computed:            true,
			},
			"platform": schema.StringAttribute{
				Description:         "Fleet platform.",
				MarkdownDescription: "The platform of the fleet.",
				Computed:            true,
			},
			"max_concurrent_sessions": schema.Int32Attribute{
				Description:         "Maximum concurrent sessions.",
				MarkdownDescription: "The maximum number of concurrent streaming sessions for an elastic fleet.",
				Computed:            true,
			},
			"max_sessions_per_instance": schema.Int32Attribute{
				Description:         "Maximum sessions per instance.",
				MarkdownDescription: "The maximum number of user sessions allowed per fleet instance.",
				Computed:            true,
			},
			"usb_device_filter_strings": schema.SetAttribute{
				Description:         "USB device filter strings.",
				MarkdownDescription: "USB device filter rules configured for the fleet.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"session_script_s3_location": schema.SingleNestedAttribute{
				Description:         "Session script S3 location.",
				MarkdownDescription: "The S3 location of the session scripts configuration ZIP file, if set.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"s3_bucket": schema.StringAttribute{
						Description:         "S3 bucket name.",
						MarkdownDescription: "The name of the S3 bucket.",
						Computed:            true,
					},
					"s3_key": schema.StringAttribute{
						Description:         "S3 object key.",
						MarkdownDescription: "The S3 object key of the session script.",
						Computed:            true,
					},
				},
			},
			"root_volume_config": schema.SingleNestedAttribute{
				Description:         "Root volume configuration.",
				MarkdownDescription: "The root volume configuration of the fleet, if set.",
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
				Description:         "Tags applied to the AppStream fleet.",
				MarkdownDescription: "Tags assigned to the AppStream fleet.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream fleet.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream fleet.",
				Computed:            true,
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the fleet was created.",
				MarkdownDescription: "The timestamp when the fleet was created, in RFC 3339 format.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				Description:         "State of the AppStream fleet.",
				MarkdownDescription: "The state of the AppStream fleet.",
				Computed:            true,
			},
			"fleet_errors": schema.SetNestedAttribute{
				Description:         "Errors reported by AWS for the fleet.",
				MarkdownDescription: "Informational list of errors reported by AWS for the fleet.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							Description:         "Error code reported by AWS.",
							MarkdownDescription: "The error code reported by AWS for the fleet.",
							Computed:            true,
						},
						"error_message": schema.StringAttribute{
							Description:         "Error message reported by AWS.",
							MarkdownDescription: "The human-readable error message reported by AWS.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
