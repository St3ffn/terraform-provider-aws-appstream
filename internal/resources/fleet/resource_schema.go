// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package fleet

import (
	"context"
	"regexp"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream Fleet",
		MarkdownDescription: "Manages an AppStream fleet. " +
			"A fleet defines the compute capacity, networking configuration, and runtime behavior " +
			"for streaming instances that host AppStream user sessions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream Fleet.",
				MarkdownDescription: "A synthetic identifier for the fleet, equal to the fleet name. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the AppStream Fleet.",
				MarkdownDescription: "The name of the AppStream fleet. " +
					"Changing this value forces the fleet to be replaced.",
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
				Description: "Name of the AppStream Image.",
				MarkdownDescription: "The name of the AppStream image used to create the fleet. " +
					"This attribute is supported only for non-elastic fleets. " +
					"Either `image_name` or `image_arn` must be specified.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"image_arn": schema.StringAttribute{
				Description: "ARN of the AppStream Image.",
				MarkdownDescription: "The ARN of the AppStream image used to create the fleet. " +
					"This attribute is supported only for non-elastic fleets. " +
					"Either `image_name` or `image_arn` must be specified.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					util.ValidARNWithServiceAndResource("appstream", "image/"),
				},
			},
			"instance_type": schema.StringAttribute{
				Description: "EC2 instance type for Fleet instances.",
				MarkdownDescription: "The EC2 instance type used by the fleet. " +
					"This field is required for non-elastic fleets.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"fleet_type": schema.StringAttribute{
				Description:         "Type of the AppStream Fleet.",
				MarkdownDescription: "The fleet type. Valid values are `ON_DEMAND`, `ALWAYS_ON`, or `ELASTIC`.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						util.AWSEnumToSlice(awstypes.FleetType.Values)...,
					),
				},
			},
			"compute_capacity": schema.SingleNestedAttribute{
				Description: "Compute capacity configuration.",
				MarkdownDescription: "Specifies the desired capacity for the fleet. Exactly one of " +
					"`desired_instances` or `desired_sessions` must be specified for non-elastic fleets. " +
					"These attributes are mutually exclusive.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"desired_instances": schema.Int32Attribute{
						Description:         "Desired number of Fleet instances.",
						MarkdownDescription: "The desired number of streaming instances for a single-session fleet.",
						Optional:            true,
						Validators: []validator.Int32{
							int32validator.AtLeast(0),
						},
					},
					"desired_sessions": schema.Int32Attribute{
						Description:         "Desired number of streaming sessions.",
						MarkdownDescription: "The desired number of concurrent streaming sessions for a multi-session fleet.",
						Optional:            true,
						Validators: []validator.Int32{
							int32validator.AtLeast(0),
						},
					},
				},
			},
			"vpc_config": schema.SingleNestedAttribute{
				Description: "VPC configuration for the Fleet.",
				MarkdownDescription: "The VPC configuration used by the fleet. " +
					"This block is required for elastic fleets.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"subnet_ids": schema.SetAttribute{
						Description: "Subnet IDs.",
						MarkdownDescription: "The subnet IDs for the fleet. This attribute is required for elastic fleets. " +
							"Elastic fleets require at least two subnets in different Availability Zones. " +
							"The Availability Zone requirement is enforced by AWS.",
						Required:    true,
						ElementType: types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"security_group_ids": schema.SetAttribute{
						Description:         "Security group IDs.",
						MarkdownDescription: "The security group IDs associated with the fleet.",
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
			"max_user_duration_in_seconds": schema.Int32Attribute{
				Description:         "Maximum user session duration.",
				MarkdownDescription: "The maximum length of time that a streaming session can remain active.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int32{
					int32validator.Between(600, 432000),
				},
			},
			"disconnect_timeout_in_seconds": schema.Int32Attribute{
				Description:         "Session disconnect timeout.",
				MarkdownDescription: "The amount of time that a disconnected session is allowed to remain active.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int32{
					int32validator.Between(60, 36000),
				},
			},
			"idle_disconnect_timeout_in_seconds": schema.Int32Attribute{
				Description: "Idle session disconnect timeout.",
				MarkdownDescription: "The amount of time, in seconds, that a session can remain idle before being disconnected. " +
					"Specify `0` to disable idle disconnection. Otherwise, the value must be a multiple of 60 seconds " +
					"between 60 and 36000 to avoid AWS rounding behavior.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int32{
					util.DurationWithStep(60, 36000, 60, true),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description of the AppStream Fleet.",
				MarkdownDescription: "The fleet description, if set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"display_name": schema.StringAttribute{
				Description:         "Display name of the AppStream Fleet.",
				MarkdownDescription: "The name displayed to users in the AppStream user interface.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"disable_imds_v1": schema.BoolAttribute{
				Description: "Disable IMDSv1 and enforce IMDSv2 for the AppStream Fleet.",
				MarkdownDescription: "Whether Instance Metadata Service Version 1 (IMDSv1) is disabled for the " +
					"AppStream fleet. If `true`, only IMDSv2 is enabled. " +
					"If `false`, both IMDSv1 and IMDSv2 are enabled.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_default_internet_access": schema.BoolAttribute{
				Description:         "Enable default internet access.",
				MarkdownDescription: "Whether instances in the fleet have access to the internet.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_join_info": schema.SingleNestedAttribute{
				Description: "Active Directory domain join configuration.",
				MarkdownDescription: "Specifies the Active Directory domain and organizational unit used to join " +
					"fleet instances to a Microsoft Active Directory domain. This configuration is not supported " +
					"for elastic fleets.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"directory_name": schema.StringAttribute{
						Description:         "Directory name.",
						MarkdownDescription: "The fully qualified domain name of the Active Directory.",
						Required:            true,
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
			"iam_role_arn": schema.StringAttribute{
				Description:         "IAM role ARN.",
				MarkdownDescription: "The ARN of the IAM role applied to fleet instances.",
				Optional:            true,
				Validators: []validator.String{
					util.ValidARNWithServiceAndResource("iam", "role/"),
				},
			},
			"stream_view": schema.StringAttribute{
				Description:         "Streaming view configuration.",
				MarkdownDescription: "Controls which streaming protocol views are enabled.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						util.AWSEnumToSlice(awstypes.StreamView.Values)...,
					),
				},
			},
			"platform": schema.StringAttribute{
				Description: "Fleet platform.",
				MarkdownDescription: "The platform of the fleet. This attribute is optional but required " +
					"for elastic fleets. If not specified, the platform is inferred from the image.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"max_concurrent_sessions": schema.Int32Attribute{
				Description: "Maximum concurrent sessions.",
				MarkdownDescription: "The maximum number of concurrent streaming sessions for an elastic fleet. " +
					"This setting is required for elastic fleets and is not allowed for other fleet types.",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"max_sessions_per_instance": schema.Int32Attribute{
				Description: "Maximum sessions per instance.",
				MarkdownDescription: "The maximum number of user sessions allowed per fleet instance. " +
					"This setting applies only to multi-session fleets.",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"usb_device_filter_strings": schema.SetAttribute{
				Description: "USB device filter strings.",
				MarkdownDescription: "Defines which USB devices can be redirected to streaming sessions when using the Windows native client. " +
					"This setting is supported only for Windows fleets. For non-Windows platforms or non-native clients, " +
					"this configuration is accepted by AWS but ignored at runtime.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeBetween(0, 100),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^((\w*)\s*(\w*)\s*,\s*(\w*)\s*,\s*\*?(\w*)\s*,\s*\*?(\w*)\s*,\s*\*?\d*\s*,\s*\*?\d*\s*,\s*[0-1]\s*,\s*[0-1]\s*)$`),
							"must match ^((\\w*)\\s*(\\w*)\\s*,\\s*(\\w*)\\s*,\\s*\\*?(\\w*)\\s*,\\s*\\*?(\\w*)\\s*,\\s*\\*?\\d*\\s*,\\s*\\*?\\d*\\s*,\\s*[0-1]\\s*,\\s*[0-1]\\s*)$",
						),
					),
				},
			},
			"session_script_s3_location": schema.SingleNestedAttribute{
				Description: "Session script S3 location.",
				MarkdownDescription: "Specifies the S3 location of the session scripts configuration ZIP file. " +
					"This setting applies only to elastic fleets.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"s3_bucket": schema.StringAttribute{
						Description:         "S3 bucket name.",
						MarkdownDescription: "The name of the S3 bucket.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(3, 63),
						},
					},
					"s3_key": schema.StringAttribute{
						Description:         "S3 object key.",
						MarkdownDescription: "The S3 object key of the session script.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 1024),
						},
					},
				},
			},
			"root_volume_config": schema.SingleNestedAttribute{
				Description:         "Root volume configuration.",
				MarkdownDescription: "Specifies the root volume configuration for fleet instances.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"volume_size_in_gb": schema.Int32Attribute{
						Description:         "Root volume size.",
						MarkdownDescription: "The size of the root volume, in GiB.",
						Optional:            true,
						Validators: []validator.Int32{
							int32validator.Between(200, 500),
						},
					},
				},
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the Fleet.",
				MarkdownDescription: "A map of tags assigned to the fleet.",
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
			"tags_all": schema.MapAttribute{
				Description:         "All tags applied to the Fleet.",
				MarkdownDescription: "A map of tags, including default tags, assigned to the fleet.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"desired_state": schema.StringAttribute{
				Description: "Desired runtime state of the AppStream Fleet.",
				MarkdownDescription: "The runtime fleet state to enforce after create and update operations. " +
					"Valid values are `INHERIT`, `RUNNING`, and `STOPPED`. " +
					"The default value is `INHERIT`, which preserves the fleet's existing runtime state.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(desiredStateInherit.String()),
				Validators: []validator.String{
					stringvalidator.OneOf(
						desiredStateInherit.String(),
						desiredStateRunning.String(),
						desiredStateStopped.String(),
					),
				},
			},
			"update_behavior": schema.StringAttribute{
				Description: "Update behavior for the AppStream Fleet lifecycle.",
				MarkdownDescription: "Controls how updates are handled when AWS requires the fleet to be stopped. " +
					"Valid values are `AUTO_STOP_START` and `FAIL_IF_RUNNING`. " +
					"If omitted, `AUTO_STOP_START` behavior is used.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(updateBehaviorAutoStopStart.String()),
				Validators: []validator.String{
					stringvalidator.OneOf(
						updateBehaviorAutoStopStart.String(),
						updateBehaviorFailIfRunning.String(),
					),
				},
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream Fleet.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream fleet.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the Fleet was created.",
				MarkdownDescription: "The timestamp when the fleet was created, in RFC 3339 format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description:         "State of the AppStream Fleet.",
				MarkdownDescription: "The state of the AppStream fleet.",
				Computed:            true,
			},
			"fleet_errors": schema.SetNestedAttribute{
				Description: "Errors reported by AWS for the Fleet.",
				MarkdownDescription: "Informational list of errors reported by AWS for the fleet. " +
					"These errors do not affect Terraform lifecycle behavior.",
				Computed: true,
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
