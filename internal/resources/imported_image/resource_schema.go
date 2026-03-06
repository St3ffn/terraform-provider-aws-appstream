// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package imported_image

import (
	"context"
	"regexp"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Schema returns the static Terraform schema for this component.
// It defines supported attributes/blocks plus validators and plan behavior.
func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream Imported Image",
		MarkdownDescription: "Manages an AppStream imported image created from an EC2 AMI. " +
			"Imported images are immutable artifacts and are typically replaced when input settings change. " +
			"After `terraform import`, configure create-time-only attributes explicitly (`iam_role_arn`, `source_ami_id`, " +
			"`agent_software_version`, `runtime_validation_config`, `app_catalog_config`) because AWS `DescribeImages` " +
			"does not return these original input values.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream Imported Image.",
				MarkdownDescription: "A synthetic identifier for the imported image, equal to the image name. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the AppStream Imported Image.",
				MarkdownDescription: "A unique name for the imported image. " +
					"Changing this value forces the imported image to be replaced.",
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
			"iam_role_arn": schema.StringAttribute{
				Description: "IAM role ARN for the Imported Image.",
				MarkdownDescription: "The ARN of the IAM role that allows AppStream to access and validate the source AMI. " +
					"Changing this value forces the imported image to be replaced. " +
					"This create-time input is not returned by `DescribeImages`, so after import it must be set explicitly in configuration.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					util.ValidARNWithServiceAndResource("iam", "role/"),
				},
			},
			"source_ami_id": schema.StringAttribute{
				Description: "Source EC2 AMI ID for the Imported Image.",
				MarkdownDescription: "The EC2 AMI ID to import into AppStream. " +
					"Changing this value forces the imported image to be replaced. " +
					"This create-time input is not returned by `DescribeImages`, so after import it must be set explicitly in configuration.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^ami-[a-z0-9]{8,17}$`),
						"must match ^ami-[a-z0-9]{8,17}$",
					),
				},
			},
			"agent_software_version": schema.StringAttribute{
				Description: "AppStream Agent software version.",
				MarkdownDescription: "The AppStream agent software version used for the imported image. " +
					"Valid values are `CURRENT_LATEST` and `ALWAYS_LATEST`. " +
					"Changing this value forces the imported image to be replaced. " +
					"This create-time input is not returned by `DescribeImages`, so after import it must be set explicitly in configuration.",
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						util.AWSEnumToSlice(awstypes.AgentSoftwareVersion.Values)...,
					),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "Display name of the AppStream Imported Image.",
				MarkdownDescription: "An optional display name for the imported image. " +
					"Changing this value forces the imported image to be replaced.",
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the AppStream Imported Image.",
				MarkdownDescription: "An optional description for the imported image. " +
					"Changing this value forces the imported image to be replaced.",
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"runtime_validation_config": schema.SingleNestedAttribute{
				Description: "Runtime validation configuration.",
				MarkdownDescription: "Configuration used to validate streaming behavior of the imported image. " +
					"Changing this value forces the imported image to be replaced. " +
					"This create-time block is not returned by `DescribeImages`, so after import it must be set explicitly in configuration.",
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"intended_instance_type": schema.StringAttribute{
						Description:         "Intended instance type for runtime validation.",
						MarkdownDescription: "The instance type used for runtime validation during image import.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
			"app_catalog_config": schema.SetNestedAttribute{
				Description: "Application catalog configuration.",
				MarkdownDescription: "Application catalog entries for the imported image. " +
					"Changing this value forces the imported image to be replaced. " +
					"This create-time block is not returned by `DescribeImages`, so after import it must be set explicitly in configuration.",
				Optional: true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(50),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Name of the AppStream Application.",
							MarkdownDescription: "The unique name of the application in the app catalog.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,99}$`),
									"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,99}$",
								),
							},
						},
						"absolute_app_path": schema.StringAttribute{
							Description:         "Absolute executable path of the Application.",
							MarkdownDescription: "The absolute path to the executable used to launch the application.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32767),
							},
						},
						"display_name": schema.StringAttribute{
							Description:         "Display name of the Application.",
							MarkdownDescription: "The display name of the application, if set.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(100),
							},
						},
						"launch_parameters": schema.StringAttribute{
							Description:         "Launch parameters of the Application.",
							MarkdownDescription: "The launch parameters passed to the application executable.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(1024),
							},
						},
						"working_directory": schema.StringAttribute{
							Description:         "Working directory of the Application.",
							MarkdownDescription: "The working directory used to launch the application.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(32767),
							},
						},
						"absolute_icon_path": schema.StringAttribute{
							Description:         "Absolute icon path of the Application.",
							MarkdownDescription: "The absolute path to the icon file for the application.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(32767),
							},
						},
						"absolute_manifest_path": schema.StringAttribute{
							Description:         "Absolute prewarm manifest path of the Application.",
							MarkdownDescription: "The absolute path to the prewarm manifest file for the application.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(32767),
							},
						},
					},
				},
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the AppStream Imported Image.",
				MarkdownDescription: "A map of tags assigned to the imported image.",
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
				Description:         "All tags applied to the AppStream Imported Image.",
				MarkdownDescription: "A map of tags, including default tags, assigned to the imported image.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream Imported Image.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream imported image.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description:         "State of the AppStream Imported Image.",
				MarkdownDescription: "The current state of the AppStream imported image.",
				Computed:            true,
			},
			"state_change_reason": schema.SingleNestedAttribute{
				Description:         "State change reason.",
				MarkdownDescription: "The reason for the most recent AppStream imported image state change, if applicable.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"code": schema.StringAttribute{
						Description:         "State change reason code.",
						MarkdownDescription: "The code describing why the AppStream imported image state changed.",
						Computed:            true,
					},
					"message": schema.StringAttribute{
						Description:         "State change reason message.",
						MarkdownDescription: "The human-readable message describing the AppStream imported image state change.",
						Computed:            true,
					},
				},
			},
			"created_time": schema.StringAttribute{
				Description:         "Creation time of the AppStream Imported Image.",
				MarkdownDescription: "The timestamp when the AppStream imported image was created, in RFC 3339 format.",
				Computed:            true,
			},
			"platform": schema.StringAttribute{
				Description:         "Platform of the AppStream Imported Image.",
				MarkdownDescription: "The operating system platform of the AppStream imported image.",
				Computed:            true,
			},
			"visibility": schema.StringAttribute{
				Description:         "Visibility of the AppStream Imported Image.",
				MarkdownDescription: "The visibility of the AppStream imported image, such as `PRIVATE` or `SHARED`.",
				Computed:            true,
			},
			"base_image_arn": schema.StringAttribute{
				Description:         "Base image ARN.",
				MarkdownDescription: "The ARN of the image from which this image was created.",
				Computed:            true,
			},
			"image_builder_supported": schema.BoolAttribute{
				Description:         "Image Builder support.",
				MarkdownDescription: "Whether an AppStream image builder can be launched from this imported image.",
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
				MarkdownDescription: "The AppStream agent version currently associated with the imported image.",
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
