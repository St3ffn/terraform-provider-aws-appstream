// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (ds *applicationDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Read an AWS AppStream Application",
		MarkdownDescription: "Reads an AppStream application. " +
			"This data source can be used to reference an existing AppStream application " +
			"that is managed outside of Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the AppStream application.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream application.",
				Computed:            true,
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream application.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream application.",
				Required:            true,
				Validators:          []validator.String{validARN()},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the AppStream application.",
				MarkdownDescription: "The name of the AppStream application.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				Description:         "Display name of the application.",
				MarkdownDescription: "The name displayed to users in the AppStream application catalog, if set.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the application.",
				MarkdownDescription: "The application description, if set.",
				Computed:            true,
			},
			"icon_s3_location": schema.SingleNestedAttribute{
				Description:         "Application icon S3 location.",
				MarkdownDescription: "Specifies the S3 location of the application icon, if set.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"s3_bucket": schema.StringAttribute{
						Description:         "S3 bucket name.",
						MarkdownDescription: "The name of the S3 bucket containing the application icon.",
						Computed:            true,
					},
					"s3_key": schema.StringAttribute{
						Description:         "S3 object key.",
						MarkdownDescription: "The S3 object key of the application icon.",
						Computed:            true,
					},
				},
			},
			"launch_path": schema.StringAttribute{
				Description:         "Application launch path.",
				MarkdownDescription: "The path to the application executable within the image.",
				Computed:            true,
			},
			"working_directory": schema.StringAttribute{
				Description:         "Application working directory.",
				MarkdownDescription: "The working directory of the application, if set.",
				Computed:            true,
			},
			"launch_parameters": schema.StringAttribute{
				Description:         "Application launch parameters.",
				MarkdownDescription: "The parameters passed to the application at launch, if set.",
				Computed:            true,
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
			"app_block_arn": schema.StringAttribute{
				Description:         "App block ARN.",
				MarkdownDescription: "The ARN of the app block associated with the application.",
				Computed:            true,
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the AppStream application.",
				MarkdownDescription: "Tags assigned to the AppStream application.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the application was created.",
				MarkdownDescription: "The timestamp when the application was created, in RFC 3339 format.",
				Computed:            true,
			},
		},
	}
}
