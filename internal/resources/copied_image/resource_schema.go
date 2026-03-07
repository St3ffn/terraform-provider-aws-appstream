// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package copied_image

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Schema returns the static Terraform schema for this component.
// It defines supported attributes/blocks plus validators and plan behavior.
func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream Copied Image",
		MarkdownDescription: "Manages an AppStream image created by copying an existing image to another region. " +
			"`CopyImage` is invoked using the provider-configured region (source region), while `destination_region` " +
			"controls where the copied image is created. " +
			"Copied images are immutable artifacts and are typically replaced when copy input settings change. " +
			"AppStream `CopyImage` does not copy tags, so tags are managed separately by the provider.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream Copied Image.",
				MarkdownDescription: "A synthetic identifier for the copied image, composed of the " +
					"destination image name, destination region, and source image name. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"destination_image_name": schema.StringAttribute{
				Description: "Name of the AppStream Copied Image.",
				MarkdownDescription: "A unique destination name for the copied image. " +
					"Changing this value forces the copied image to be replaced.",
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
			"destination_image_description": schema.StringAttribute{
				Description: "Description of the AppStream Copied Image.",
				MarkdownDescription: "A description for the copied image. " +
					"Changing this value forces the copied image to be replaced.",
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"destination_region": schema.StringAttribute{
				Description: "Destination AWS region for the copy operation.",
				MarkdownDescription: "The AWS region where the copied image is created. " +
					"Changing this value forces the copied image to be replaced.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
				},
			},
			"source_image_name": schema.StringAttribute{
				Description: "Source AppStream image name.",
				MarkdownDescription: "The name of the source image used for the copy operation. " +
					"Changing this value forces the copied image to be replaced. " +
					"AWS read APIs do not reliably return this value, so import requires it as the third segment " +
					"of the import identifier.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
						"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
					)},
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the AppStream Copied Image.",
				MarkdownDescription: "A map of tags assigned to the copied image.",
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
				Description:         "All tags applied to the AppStream Copied Image.",
				MarkdownDescription: "A map of tags, including default tags, assigned to the copied image.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream Copied Image.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream copied image.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description:         "State of the AppStream Copied Image.",
				MarkdownDescription: "The current state of the AppStream copied image.",
				Computed:            true,
			},
			"state_change_reason": schema.SingleNestedAttribute{
				Description:         "State change reason.",
				MarkdownDescription: "The reason for the most recent AppStream copied image state change, if applicable.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"code": schema.StringAttribute{
						Description:         "State change reason code.",
						MarkdownDescription: "The code describing why the AppStream copied image state changed.",
						Computed:            true,
					},
					"message": schema.StringAttribute{
						Description:         "State change reason message.",
						MarkdownDescription: "The human-readable message describing the AppStream copied image state change.",
						Computed:            true,
					},
				},
			},
			"created_time": schema.StringAttribute{
				Description:         "Creation time of the AppStream Copied Image.",
				MarkdownDescription: "The timestamp when the AppStream copied image was created, in RFC 3339 format.",
				Computed:            true,
			},
			"platform": schema.StringAttribute{
				Description:         "Platform of the AppStream Copied Image.",
				MarkdownDescription: "The operating system platform of the AppStream copied image.",
				Computed:            true,
			},
			"visibility": schema.StringAttribute{
				Description:         "Visibility of the AppStream Copied Image.",
				MarkdownDescription: "The visibility of the AppStream copied image, such as `PRIVATE` or `SHARED`.",
				Computed:            true,
			},
			"base_image_arn": schema.StringAttribute{
				Description:         "Base image ARN.",
				MarkdownDescription: "The ARN of the image from which this image was created.",
				Computed:            true,
			},
			"image_builder_supported": schema.BoolAttribute{
				Description:         "Image Builder support.",
				MarkdownDescription: "Whether an AppStream image builder can be launched from this copied image.",
				Computed:            true,
			},
			"image_builder_name": schema.StringAttribute{
				Description:         "Image Builder name.",
				MarkdownDescription: "The name of the image builder used to create the image, if applicable.",
				Computed:            true,
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
			"appstream_agent_version": schema.StringAttribute{
				Description:         "AppStream agent version.",
				MarkdownDescription: "The AppStream agent version currently associated with the copied image.",
				Computed:            true,
			},
			"public_base_image_released_date": schema.StringAttribute{
				Description:         "Public base image release date.",
				MarkdownDescription: "The release date of the public base image, in RFC 3339 format, if available.",
				Computed:            true,
			},
			"image_permissions": schema.SingleNestedAttribute{
				Description:         "Image permissions.",
				MarkdownDescription: "Permissions granted for the image.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"allow_fleet": schema.BoolAttribute{
						Description:         "Allow Fleet usage.",
						MarkdownDescription: "Whether the image can be used by fleets.",
						Computed:            true,
					},
					"allow_image_builder": schema.BoolAttribute{
						Description:         "Allow Image Builder usage.",
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
		},
	}
}
