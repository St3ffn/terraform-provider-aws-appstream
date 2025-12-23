// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *stackResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream Stack",
		MarkdownDescription: "Manages an AppStream stack. " +
			"A stack defines the streaming configuration and user experience for AppStream sessions, " +
			"including storage, networking, user actions, and application settings persistence.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream stack.",
				MarkdownDescription: "A synthetic identifier for the stack, equal to the stack name. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the AppStream stack.",
				MarkdownDescription: "The name of the AppStream stack. " +
					"Changing this value forces the stack to be replaced.",
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
			"description": schema.StringAttribute{
				Description: "Description of the AppStream stack.",
				MarkdownDescription: "The stack description, if set. " +
					"Must be 256 characters or fewer.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"display_name": schema.StringAttribute{
				Description:         "Display name of the AppStream stack.",
				MarkdownDescription: "The name displayed to users in the AppStream user interface.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"storage_connectors": schema.SetNestedAttribute{
				Description:         "Storage connectors for the stack.",
				MarkdownDescription: "Storage connectors that enable persistent storage for users of the stack.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"connector_type": schema.StringAttribute{
							Description:         "Type of the storage connector.",
							MarkdownDescription: "The type of storage connector to enable for the stack.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"HOMEFOLDERS",
									"GOOGLE_DRIVE",
									"ONE_DRIVE",
								),
							},
						},
						"resource_identifier": schema.StringAttribute{
							Description:         "Resource identifier of the storage connector.",
							MarkdownDescription: "The resource identifier associated with the storage connector.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 2048),
							},
						},
						"domains": schema.SetAttribute{
							Description:         "Domains associated with the storage connector.",
							MarkdownDescription: "The domain names associated with the storage connector.",
							Optional:            true,
							ElementType:         types.StringType,
							Validators: []validator.Set{
								setvalidator.SizeAtMost(50),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 64),
								),
							},
						},
						"domains_require_admin_consent": schema.SetAttribute{
							Description:         "Domains requiring administrator consent.",
							MarkdownDescription: "The domain names that require administrator consent before access is granted.",
							Optional:            true,
							ElementType:         types.StringType,
							Validators: []validator.Set{
								setvalidator.SizeAtMost(50),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 64),
								),
							},
						},
					},
				},
			},
			"redirect_url": schema.StringAttribute{
				Description:         "Redirect URL after streaming sessions end.",
				MarkdownDescription: "The URL users are redirected to after their AppStream streaming session ends.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1000),
				},
			},
			"feedback_url": schema.StringAttribute{
				Description:         "Feedback URL for the stack.",
				MarkdownDescription: "The URL users are redirected to after clicking the **Send Feedback** link.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1000),
				},
			},
			"user_settings": schema.SetNestedAttribute{
				Description:         "User settings for streaming sessions.",
				MarkdownDescription: "Actions that are enabled or disabled for users during streaming sessions.",
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Description:         "User action.",
							MarkdownDescription: "The user action that can be enabled or disabled during a streaming session.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"CLIPBOARD_COPY_FROM_LOCAL_DEVICE",
									"CLIPBOARD_COPY_TO_LOCAL_DEVICE",
									"FILE_UPLOAD",
									"FILE_DOWNLOAD",
									"PRINTING_TO_LOCAL_DEVICE",
									"DOMAIN_PASSWORD_SIGNIN",
									"DOMAIN_SMART_CARD_SIGNIN",
									"AUTO_TIME_ZONE_REDIRECTION",
								),
							},
						},
						"permission": schema.StringAttribute{
							Description:         "Permission for the action.",
							MarkdownDescription: "Specifies whether the action is enabled or disabled.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("ENABLED", "DISABLED"),
							},
						},
						"maximum_length": schema.Int32Attribute{
							Description: "Maximum number of characters that can be copied.",
							MarkdownDescription: "Specifies the maximum number of characters that can be copied for clipboard actions. " +
								"This setting applies only to `CLIPBOARD_COPY_FROM_LOCAL_DEVICE` and `CLIPBOARD_COPY_TO_LOCAL_DEVICE`. " +
								"It cannot be set when permission is `DISABLED`.",
							Optional: true,
							Validators: []validator.Int32{
								int32validator.Between(1, 20971520),
							},
						},
					},
				},
			},
			"application_settings": schema.SingleNestedAttribute{
				Description:         "Application settings persistence configuration.",
				MarkdownDescription: "Controls persistence of application settings for users of the stack.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description:         "Enable application settings persistence.",
						MarkdownDescription: "Whether application settings persistence is enabled.",
						Required:            true,
					},
					"settings_group": schema.StringAttribute{
						Description:         "Name of the application settings group.",
						MarkdownDescription: "The name of the application settings group.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(100),
						},
					},
					"s3_bucket_name": schema.StringAttribute{
						Description:         "S3 bucket name for persistent application settings.",
						MarkdownDescription: "The S3 bucket name where users persistent application settings are stored.",
						Computed:            true,
					},
				},
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the AppStream stack.",
				MarkdownDescription: "A map of tags assigned to the AppStream stack.",
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
			"access_endpoints": schema.SetNestedAttribute{
				Description:         "VPC access endpoints for the stack.",
				MarkdownDescription: "Interface VPC endpoints through which users can connect to the stack.",
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 4),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"endpoint_type": schema.StringAttribute{
							Description:         "Type of the access endpoint.",
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
			"embed_host_domains": schema.SetAttribute{
				Description:         "Domains allowed for embedded streaming.",
				MarkdownDescription: "Domains where streaming sessions can be embedded in an iframe.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 20),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtMost(128),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`),
							"must match (?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]",
						),
					),
				},
			},
			"streaming_experience_settings": schema.SingleNestedAttribute{
				Description:         "Streaming experience configuration.",
				MarkdownDescription: "Controls the preferred streaming protocol for the stack.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"preferred_protocol": schema.StringAttribute{
						Description:         "Preferred streaming protocol.",
						MarkdownDescription: "The preferred streaming protocol for the stack.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("TCP", "UDP"),
						},
					},
				},
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream stack.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream stack.",
				Computed:            true,
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the stack was created.",
				MarkdownDescription: "The timestamp when the stack was created, in RFC 3339 format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"stack_errors": schema.SetNestedAttribute{
				Description: "Errors reported by AWS for the stack.",
				MarkdownDescription: "Informational list of errors reported by AWS for the stack. " +
					"These errors do not affect Terraform lifecycle behavior.",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							Description:         "Error code reported by AWS.",
							MarkdownDescription: "The error code reported by AWS for the stack.",
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
