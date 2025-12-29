// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package application

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an AWS AppStream Application",
		MarkdownDescription: "Manages an AppStream application. " +
			"An application defines how software is launched on AppStream Elastic fleet streaming instances, " +
			"including launch configuration, supported platforms, instance families, and application metadata.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the AppStream application.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream application. " +
					"This value is managed by the provider and cannot be set manually.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the AppStream application.",
				MarkdownDescription: "The name of the AppStream application. " +
					"Changing this value forces the application to be replaced.",
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
				Description:         "Display name of the application.",
				MarkdownDescription: "The name displayed to users in the AppStream application catalog.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description of the application.",
				MarkdownDescription: "The application description, if set. Must be 256 characters or fewer.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"icon_s3_location": schema.SingleNestedAttribute{
				Description: "Application icon S3 location.",
				MarkdownDescription: "Specifies the S3 location of the application icon. " +
					"The icon is displayed to users in the AppStream application catalog.",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"s3_bucket": schema.StringAttribute{
						Description:         "S3 bucket name.",
						MarkdownDescription: "The name of the S3 bucket containing the application icon.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(3, 63),
						},
					},
					"s3_key": schema.StringAttribute{
						Description:         "S3 object key.",
						MarkdownDescription: "The S3 object key of the application icon.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 1024),
						},
					},
				},
			},
			"launch_path": schema.StringAttribute{
				Description:         "Application launch path.",
				MarkdownDescription: "The path to the application executable within the image.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"working_directory": schema.StringAttribute{
				Description:         "Application working directory.",
				MarkdownDescription: "The working directory of the application.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"launch_parameters": schema.StringAttribute{
				Description:         "Application launch parameters.",
				MarkdownDescription: "The parameters passed to the application at launch.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"platforms": schema.SetAttribute{
				Description:         "Supported platforms.",
				MarkdownDescription: "The platforms on which the application can run.",
				Required:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(4),
				},
			},
			"instance_families": schema.SetAttribute{
				Description:         "Supported instance families.",
				MarkdownDescription: "The instance families supported by the application.",
				Required:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"app_block_arn": schema.StringAttribute{
				Description:         "App block ARN.",
				MarkdownDescription: "The ARN of the app block associated with the application.",
				Required:            true,
				Validators: []validator.String{
					util.ValidARNWithServiceAndResource("appstream", "app-block/"),
				},
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the AppStream application.",
				MarkdownDescription: "A map of tags assigned to the AppStream application.",
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
				Description:         "ARN of the AppStream application.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream application.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the application was created.",
				MarkdownDescription: "The timestamp when the application was created, in RFC 3339 format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
