// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream App Block",
		MarkdownDescription: "Manages an AppStream app block. " +
			"An app block defines a reusable application package for AppStream Elastic fleets, " +
			"including application binaries and setup scripts stored in Amazon S3. " +
			"App blocks can be referenced by multiple AppStream applications.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream app block.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream app block. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the AppStream app block.",
				MarkdownDescription: "The name of the AppStream app block. " +
					"Changing this value forces the app block to be replaced.",
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
			"display_name": schema.StringAttribute{
				Description:         "Display name of the app block.",
				MarkdownDescription: "The display name of the AppStream app block.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description of the app block.",
				MarkdownDescription: "The AppStream app block description, if set. Must be 256 characters or fewer.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"source_s3_location": schema.SingleNestedAttribute{
				Description: "Source S3 location of the app block.",
				MarkdownDescription: "Specifies the Amazon S3 location that contains the app block source. " +
					"Depending on the packaging type, this can be an application package (VHD) or " +
					"source content used to build the app block.",
				Required: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"s3_bucket": schema.StringAttribute{
						Description:         "S3 bucket name.",
						MarkdownDescription: "The name of the Amazon S3 bucket.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(3, 63),
						},
					},
					"s3_key": schema.StringAttribute{
						Description:         "S3 object key.",
						MarkdownDescription: "The S3 object key of the app block source.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 1024),
						},
					},
				},
			},
			"setup_script_details": schema.SingleNestedAttribute{
				Description: "Setup script configuration.",
				MarkdownDescription: "Specifies the setup script that is executed when the app block is built. " +
					"This configuration is required for app blocks with the `CUSTOM` packaging type.",
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"script_s3_location": schema.SingleNestedAttribute{
						Description:         "Setup script S3 location.",
						MarkdownDescription: "Specifies the Amazon S3 location of the setup script.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"s3_bucket": schema.StringAttribute{
								Description:         "S3 bucket name.",
								MarkdownDescription: "The name of the Amazon S3 bucket.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.LengthBetween(3, 63),
								},
							},
							"s3_key": schema.StringAttribute{
								Description:         "S3 object key.",
								MarkdownDescription: "The S3 object key of the setup script.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 1024),
								},
							},
						},
					},
					"executable_path": schema.StringAttribute{
						Description:         "Run path for the script.",
						MarkdownDescription: "The run path for the setup script.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"executable_parameters": schema.StringAttribute{
						Description:         "Runtime parameters for the script.",
						MarkdownDescription: "The runtime parameters passed to the setup script.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"timeout_in_seconds": schema.Int32Attribute{
						Description:         "Setup script timeout.",
						MarkdownDescription: "The maximum amount of time, in seconds, that the setup script is allowed to run.",
						Required:            true,
						Validators: []validator.Int32{
							int32validator.AtLeast(0),
						},
					},
				},
			},
			"post_setup_script_details": schema.SingleNestedAttribute{
				Description: "Post setup script configuration.",
				MarkdownDescription: "Specifies a post-setup script that is executed after the app block is created. " +
					"This configuration is supported only for app blocks with the `APPSTREAM2` packaging type.",
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"script_s3_location": schema.SingleNestedAttribute{
						Description:         "Post setup script S3 location.",
						MarkdownDescription: "Specifies the Amazon S3 location of the post-setup script.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"s3_bucket": schema.StringAttribute{
								Description:         "S3 bucket name.",
								MarkdownDescription: "The name of the Amazon S3 bucket.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.LengthBetween(3, 63),
								},
							},
							"s3_key": schema.StringAttribute{
								Description:         "S3 object key.",
								MarkdownDescription: "The S3 object key of the post-setup script.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 1024),
								},
							},
						},
					},
					"executable_path": schema.StringAttribute{
						Description:         "Run path for the script.",
						MarkdownDescription: "The run path for the post-setup script.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"executable_parameters": schema.StringAttribute{
						Description:         "Runtime parameters for the script.",
						MarkdownDescription: "The runtime parameters passed to the post-setup script.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"timeout_in_seconds": schema.Int32Attribute{
						Description:         "Post setup script timeout.",
						MarkdownDescription: "The maximum amount of time, in seconds, that the post-setup script is allowed to run.",
						Required:            true,
						Validators: []validator.Int32{
							int32validator.AtLeast(0),
						},
					},
				},
			},
			"packaging_type": schema.StringAttribute{
				Description:         "Packaging type of the app block.",
				MarkdownDescription: "The packaging type of the app block. Valid values are `CUSTOM` or `APPSTREAM2`.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("CUSTOM", "APPSTREAM2"),
				},
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the app block.",
				MarkdownDescription: "A map of tags assigned to the AppStream app block.",
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
				Description:         "ARN of the AppStream app block.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream app block.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the app block was created.",
				MarkdownDescription: "The timestamp when the app block was created, in RFC 3339 format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_block_errors": schema.SetNestedAttribute{
				Description: "Errors reported by AWS for the app block.",
				MarkdownDescription: "Informational list of errors reported by AWS for the app block. " +
					"These errors do not affect Terraform lifecycle behavior.",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							Description:         "Error code reported by AWS.",
							MarkdownDescription: "The error code reported by AWS for the app block.",
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
