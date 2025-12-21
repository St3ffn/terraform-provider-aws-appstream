// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (ds *stackDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read an AWS AppStream Stack",
		MarkdownDescription: "Reads an AppStream stack. " +
			"This data source can be used to reference an existing AppStream stack that is managed outside of Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the AppStream stack.",
				MarkdownDescription: "A synthetic identifier for the stack, equal to the stack name.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the AppStream stack.",
				MarkdownDescription: "The name of the AppStream stack to read.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description of the AppStream stack.",
				MarkdownDescription: "The stack description, if set.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				Description:         "Display name of the AppStream stack.",
				MarkdownDescription: "The name displayed to users in the AppStream user interface.",
				Computed:            true,
			},
			"storage_connectors": schema.SetNestedAttribute{
				Description:         "Storage connectors for the stack.",
				MarkdownDescription: "Storage connectors that are enabled for the stack.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"connector_type": schema.StringAttribute{
							Description:         "Type of the storage connector.",
							MarkdownDescription: "The type of storage connector.",
							Computed:            true,
						},
						"resource_identifier": schema.StringAttribute{
							Description:         "Resource identifier of the storage connector.",
							MarkdownDescription: "The resource identifier associated with the storage connector.",
							Computed:            true,
						},
						"domains": schema.SetAttribute{
							Description:         "Domains associated with the storage connector.",
							MarkdownDescription: "The domain names associated with the storage connector.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"domains_require_admin_consent": schema.SetAttribute{
							Description:         "Domains requiring administrator consent.",
							MarkdownDescription: "The domain names that require administrator consent.",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
			"redirect_url": schema.StringAttribute{
				Description:         "Redirect URL after streaming sessions end.",
				MarkdownDescription: "The URL users are redirected to after their AppStream streaming session ends.",
				Computed:            true,
			},
			"feedback_url": schema.StringAttribute{
				Description:         "Feedback URL for the stack.",
				MarkdownDescription: "The URL users are redirected to after clicking the **Send Feedback** link.",
				Computed:            true,
			},
			"user_settings": schema.SetNestedAttribute{
				Description:         "User settings for streaming sessions.",
				MarkdownDescription: "Actions that are enabled or disabled for users during streaming sessions.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Description:         "User action.",
							MarkdownDescription: "The user action.",
							Computed:            true,
						},
						"permission": schema.StringAttribute{
							Description:         "Permission for the action.",
							MarkdownDescription: "Whether the action is enabled or disabled.",
							Computed:            true,
						},
						"maximum_length": schema.Int32Attribute{
							Description:         "Maximum number of characters that can be copied.",
							MarkdownDescription: "The maximum number of characters that can be copied for clipboard actions.",
							Computed:            true,
						},
					},
				},
			},
			"application_settings": schema.SingleNestedAttribute{
				Description:         "Application settings persistence configuration.",
				MarkdownDescription: "Application settings persistence configuration for the stack.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description:         "Enable application settings persistence.",
						MarkdownDescription: "Whether application settings persistence is enabled.",
						Computed:            true,
					},
					"settings_group": schema.StringAttribute{
						Description:         "Name of the application settings group.",
						MarkdownDescription: "The name of the application settings group.",
						Computed:            true,
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
				MarkdownDescription: "Tags assigned to the AppStream stack.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"access_endpoints": schema.SetNestedAttribute{
				Description:         "VPC access endpoints for the stack.",
				MarkdownDescription: "Interface VPC endpoints through which users can connect to the stack.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"endpoint_type": schema.StringAttribute{
							Description:         "Type of the access endpoint.",
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
			"embed_host_domains": schema.SetAttribute{
				Description:         "Domains allowed for embedded streaming.",
				MarkdownDescription: "Domains where streaming sessions can be embedded in an iframe.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"streaming_experience_settings": schema.SingleNestedAttribute{
				Description:         "Streaming experience configuration.",
				MarkdownDescription: "Preferred streaming protocol configuration for the stack.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"preferred_protocol": schema.StringAttribute{
						Description:         "Preferred streaming protocol.",
						MarkdownDescription: "The preferred streaming protocol for the stack.",
						Computed:            true,
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
			},
			"stack_errors": schema.SetNestedAttribute{
				Description:         "Errors reported by AWS for the stack.",
				MarkdownDescription: "Informational list of errors reported by AWS for the stack.",
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
					},
				},
			},
		},
	}
}
