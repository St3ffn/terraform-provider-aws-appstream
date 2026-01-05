// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (ds *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read an AWS AppStream Image",
		MarkdownDescription: "Reads an AppStream image. " +
			"This data source can be used to reference an existing AppStream image that is managed outside of Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the AppStream image.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream image.",
				Computed:            true,
			},
			"arn": schema.StringAttribute{
				Description: "ARN of the AppStream image.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream image. " +
					"Cannot be used together with `name` or `name_regex`.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					util.ValidARNWithServiceAndResource("appstream", "image/"),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the AppStream image.",
				MarkdownDescription: "The name of the AppStream image. " +
					"Cannot be used together with `arn` or `name_regex`. " +
					"If multiple images with the same name exist, the data source will return an error.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					),
				},
			},
			"name_regex": schema.StringAttribute{
				Description: "Regular expression to match image names.",
				MarkdownDescription: "A regular expression used to match AppStream image names. " +
					"Uses Go regular expression syntax. " +
					"Cannot be used together with `arn` or `name`. " +
					"If the expression matches multiple images, the data source will return an error.",
				Optional: true,
				Validators: []validator.String{
					util.ValidRegex(),
				},
			},
			"visibility": schema.StringAttribute{
				Description:         "Visibility of the AppStream image.",
				MarkdownDescription: "The image visibility. Valid values are `PUBLIC`, `PRIVATE`, or `SHARED`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"PUBLIC",
						"PRIVATE",
						"SHARED",
					),
				},
			},
			"most_recent": schema.BoolAttribute{
				Description: "Return the most recent image when multiple matches exist.",
				MarkdownDescription: "Whether the most recent AppStream image is returned when multiple images match the selection criteria. " +
					"If set to `false` and multiple images match, the data source will return an error.",
				Optional: true,
			},
			"base_image_arn": schema.StringAttribute{
				Description:         "Base image ARN.",
				MarkdownDescription: "The ARN of the image from which this image was created.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				Description:         "Display name of the image.",
				MarkdownDescription: "The name displayed to users for the image, if set.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				Description:         "Image state.",
				MarkdownDescription: "The current state of the image.",
				Computed:            true,
			},
			"image_builder_supported": schema.BoolAttribute{
				Description:         "Image builder support.",
				MarkdownDescription: "Whether an image builder can be launched from this image.",
				Computed:            true,
			},
			"image_builder_name": schema.StringAttribute{
				Description:         "Image builder name.",
				MarkdownDescription: "The name of the image builder used to create the image, if applicable.",
				Computed:            true,
			},
			"platform": schema.StringAttribute{
				Description:         "Image platform.",
				MarkdownDescription: "The operating system platform of the image.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the image.",
				MarkdownDescription: "The image description, if set.",
				Computed:            true,
			},
			"state_change_reason": schema.SingleNestedAttribute{
				Description:         "State change reason.",
				MarkdownDescription: "The reason for the most recent image state change, if applicable.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"code": schema.StringAttribute{
						Description:         "State change reason code.",
						MarkdownDescription: "The code describing why the image state changed.",
						Computed:            true,
					},
					"message": schema.StringAttribute{
						Description:         "State change reason message.",
						MarkdownDescription: "The human-readable message describing the state change.",
						Computed:            true,
					},
				},
			},
			"applications": schema.SetNestedAttribute{
				Description:         "Applications included in the image.",
				MarkdownDescription: "Applications that are associated with the image.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Application name.",
							MarkdownDescription: "The name of the application.",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							Description:         "Application display name.",
							MarkdownDescription: "The display name of the application.",
							Computed:            true,
						},
						"icon_url": schema.StringAttribute{
							Description:         "Application icon URL.",
							MarkdownDescription: "The URL of the application icon, if set.",
							Computed:            true,
						},
						"launch_path": schema.StringAttribute{
							Description:         "Application launch path.",
							MarkdownDescription: "The path to the application executable.",
							Computed:            true,
						},
						"launch_parameters": schema.StringAttribute{
							Description:         "Application launch parameters.",
							MarkdownDescription: "The parameters passed to the application at launch.",
							Computed:            true,
						},
						"enabled": schema.BoolAttribute{
							Description:         "Application enabled.",
							MarkdownDescription: "Whether the application is enabled.",
							Computed:            true,
						},
						"metadata": schema.MapAttribute{
							Description:         "Application metadata.",
							MarkdownDescription: "Additional metadata associated with the application.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"working_directory": schema.StringAttribute{
							Description:         "Application working directory.",
							MarkdownDescription: "The working directory of the application.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "Application description.",
							MarkdownDescription: "The application description, if set.",
							Computed:            true,
						},
						"arn": schema.StringAttribute{
							Description:         "Application ARN.",
							MarkdownDescription: "The ARN of the application.",
							Computed:            true,
						},
						"app_block_arn": schema.StringAttribute{
							Description:         "App block ARN.",
							MarkdownDescription: "The ARN of the app block associated with the application.",
							Computed:            true,
						},
						"icon_s3_location": schema.SingleNestedAttribute{
							Description:         "Application icon S3 location.",
							MarkdownDescription: "The S3 location of the application icon, if set.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"s3_bucket": schema.StringAttribute{
									Description:         "S3 bucket name.",
									MarkdownDescription: "The name of the S3 bucket.",
									Computed:            true,
								},
								"s3_key": schema.StringAttribute{
									Description:         "S3 object key.",
									MarkdownDescription: "The S3 object key of the icon.",
									Computed:            true,
								},
							},
						},
						"platforms": schema.SetAttribute{
							Description:         "Supported platforms.",
							MarkdownDescription: "The platforms on which the application can run.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"instance_families": schema.SetAttribute{
							Description:         "Supported instance families.",
							MarkdownDescription: "The instance families supported by the application.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"created_time": schema.StringAttribute{
							Description:         "Application creation time.",
							MarkdownDescription: "The timestamp when the application was created, in RFC 3339 format.",
							Computed:            true,
						},
					},
				},
			},
			"created_time": schema.StringAttribute{
				Description:         "Image creation time.",
				MarkdownDescription: "The timestamp when the image was created, in RFC 3339 format.",
				Computed:            true,
			},
			"public_base_image_released_date": schema.StringAttribute{
				Description:         "Public base image release date.",
				MarkdownDescription: "The release date of the public base image, in RFC 3339 format.",
				Computed:            true,
			},
			"appstream_agent_version": schema.StringAttribute{
				Description:         "AppStream agent version.",
				MarkdownDescription: "The AppStream agent version used by the image.",
				Computed:            true,
			},
			"image_permissions": schema.SingleNestedAttribute{
				Description:         "Image permissions.",
				MarkdownDescription: "Permissions granted for the image.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"allow_fleet": schema.BoolAttribute{
						Description:         "Allow fleet usage.",
						MarkdownDescription: "Whether the image can be used by fleets.",
						Computed:            true,
					},
					"allow_image_builder": schema.BoolAttribute{
						Description:         "Allow image builder usage.",
						MarkdownDescription: "Whether the image can be used by image builders.",
						Computed:            true,
					},
				},
			},
			"image_errors": schema.SetNestedAttribute{
				Description:         "Image errors.",
				MarkdownDescription: "Errors reported by AWS during image creation or management.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							Description:         "Error code.",
							MarkdownDescription: "The error code reported by AWS.",
							Computed:            true,
						},
						"error_message": schema.StringAttribute{
							Description:         "Error message.",
							MarkdownDescription: "The human-readable error message.",
							Computed:            true,
						},
						"error_timestamp": schema.StringAttribute{
							Description:         "Error timestamp.",
							MarkdownDescription: "The time the error occurred, in RFC 3339 format.",
							Computed:            true,
						},
					},
				},
			},
			"latest_appstream_agent_version": schema.StringAttribute{
				Description:         "Latest AppStream agent version.",
				MarkdownDescription: "Indicates whether the image uses the latest AppStream agent version.",
				Computed:            true,
			},
			"supported_instance_families": schema.SetAttribute{
				Description:         "Supported instance families.",
				MarkdownDescription: "The instance families supported by the image.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"dynamic_app_providers_enabled": schema.StringAttribute{
				Description:         "Dynamic app providers status.",
				MarkdownDescription: "Indicates whether dynamic app providers are enabled.",
				Computed:            true,
			},
			"image_shared_with_others": schema.StringAttribute{
				Description:         "Image sharing status.",
				MarkdownDescription: "Indicates whether the image is shared with other AWS accounts.",
				Computed:            true,
			},
			"managed_software_included": schema.BoolAttribute{
				Description:         "Managed software included.",
				MarkdownDescription: "Whether the image includes managed software.",
				Computed:            true,
			},
			"image_type": schema.StringAttribute{
				Description:         "Image type.",
				MarkdownDescription: "The type of the image.",
				Computed:            true,
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the AppStream image.",
				MarkdownDescription: "Tags assigned to the AppStream image.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}
