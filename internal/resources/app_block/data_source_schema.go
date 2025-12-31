// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (ds *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read an AWS AppStream App Block",
		MarkdownDescription: "Reads an AppStream app block. " +
			"This data source can be used to reference an existing AppStream app block " +
			"that is managed outside of Terraform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the AppStream app block.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream app block.",
				Computed:            true,
			},
			"arn": schema.StringAttribute{
				Description:         "ARN of the AppStream app block.",
				MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream app block.",
				Required:            true,
				Validators: []validator.String{
					util.ValidARNWithServiceAndResource("appstream", "app-block/"),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the AppStream app block.",
				MarkdownDescription: "The name of the AppStream app block.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				Description:         "Display name of the app block.",
				MarkdownDescription: "The display name of the AppStream app block, if set.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the app block.",
				MarkdownDescription: "The app block description, if set.",
				Computed:            true,
			},
			"source_s3_location": schema.SingleNestedAttribute{
				Description:         "Source S3 location of the app block.",
				MarkdownDescription: "Specifies the Amazon S3 location that contains the app block source.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"s3_bucket": schema.StringAttribute{
						Description:         "S3 bucket name.",
						MarkdownDescription: "The name of the Amazon S3 bucket.",
						Computed:            true,
					},
					"s3_key": schema.StringAttribute{
						Description:         "S3 object key.",
						MarkdownDescription: "The S3 object key of the app block source, if set.",
						Computed:            true,
					},
				},
			},
			"setup_script_details": schema.SingleNestedAttribute{
				Description:         "Setup script configuration.",
				MarkdownDescription: "Specifies the setup script used to build the app block, if present.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"script_s3_location": schema.SingleNestedAttribute{
						Description:         "Setup script S3 location.",
						MarkdownDescription: "Specifies the Amazon S3 location of the setup script.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"s3_bucket": schema.StringAttribute{
								Description:         "S3 bucket name.",
								MarkdownDescription: "The name of the Amazon S3 bucket.",
								Computed:            true,
							},
							"s3_key": schema.StringAttribute{
								Description:         "S3 object key.",
								MarkdownDescription: "The S3 object key of the setup script.",
								Computed:            true,
							},
						},
					},
					"executable_path": schema.StringAttribute{
						Description:         "Run path for the script.",
						MarkdownDescription: "The run path for the setup script.",
						Computed:            true,
					},
					"executable_parameters": schema.StringAttribute{
						Description:         "Runtime parameters for the script.",
						MarkdownDescription: "The runtime parameters passed to the setup script, if set.",
						Computed:            true,
					},
					"timeout_in_seconds": schema.Int32Attribute{
						Description:         "Setup script timeout.",
						MarkdownDescription: "The maximum amount of time, in seconds, that the setup script is allowed to run.",
						Computed:            true,
					},
				},
			},
			"post_setup_script_details": schema.SingleNestedAttribute{
				Description:         "Post setup script configuration.",
				MarkdownDescription: "Specifies the post-setup script executed after app block creation, if present.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"script_s3_location": schema.SingleNestedAttribute{
						Description:         "Post setup script S3 location.",
						MarkdownDescription: "Specifies the Amazon S3 location of the post-setup script.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"s3_bucket": schema.StringAttribute{
								Description:         "S3 bucket name.",
								MarkdownDescription: "The name of the Amazon S3 bucket.",
								Computed:            true,
							},
							"s3_key": schema.StringAttribute{
								Description:         "S3 object key.",
								MarkdownDescription: "The S3 object key of the post-setup script.",
								Computed:            true,
							},
						},
					},
					"executable_path": schema.StringAttribute{
						Description:         "Run path for the script.",
						MarkdownDescription: "The run path for the post-setup script.",
						Computed:            true,
					},
					"executable_parameters": schema.StringAttribute{
						Description:         "Runtime parameters for the script.",
						MarkdownDescription: "The runtime parameters passed to the post-setup script, if set.",
						Computed:            true,
					},
					"timeout_in_seconds": schema.Int32Attribute{
						Description:         "Post setup script timeout.",
						MarkdownDescription: "The maximum amount of time, in seconds, that the post-setup script is allowed to run.",
						Computed:            true,
					},
				},
			},
			"packaging_type": schema.StringAttribute{
				Description:         "Packaging type of the app block.",
				MarkdownDescription: "The packaging type of the app block.",
				Computed:            true,
			},
			"tags": schema.MapAttribute{
				Description:         "Tags applied to the app block.",
				MarkdownDescription: "Tags assigned to the AppStream app block.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"created_time": schema.StringAttribute{
				Description:         "Time the app block was created.",
				MarkdownDescription: "The timestamp when the app block was created, in RFC 3339 format.",
				Computed:            true,
			},
			"app_block_errors": schema.SetNestedAttribute{
				Description:         "Errors reported by AWS for the app block.",
				MarkdownDescription: "Informational list of errors reported by AWS for the app block.",
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
