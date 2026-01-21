// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package export_image_tasks

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (ds *dataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Read AWS AppStream Export Image Tasks",
		MarkdownDescription: "Reads multiple AppStream export image tasks for the current AWS account and region. " +
			"This data source can be used to monitor the status of export operations, list completed or failed tasks, " +
			"and inspect metadata such as created time, resulting AMI IDs, and error details.",
		Attributes: map[string]schema.Attribute{
			"filters": schema.SetNestedAttribute{
				Description: "Filters for export image tasks.",
				MarkdownDescription: "Optional filters to narrow down the list of export image tasks. " +
					"Filters allow selecting tasks by attributes such as state or creation date.",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Filter name.",
							MarkdownDescription: "The name of the filter to apply.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`),
									"must match ^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$",
								),
							},
						},
						"values": schema.SetAttribute{
							Description:         "Filter values.",
							MarkdownDescription: "One or more values for the specified filter name.",
							Required:            true,
							ElementType:         types.StringType,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(
										regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_:/.-]{0,200}$`),
										"must match ^[a-zA-Z0-9][a-zA-Z0-9_:/.-]{0,200}$",
									),
								),
							},
						},
					},
				},
			},
			"export_image_tasks": schema.SetNestedAttribute{
				Description:         "Export image tasks.",
				MarkdownDescription: "The list of export image tasks that match the specified filters.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"task_id": schema.StringAttribute{
							Description:         "Export image task ID.",
							MarkdownDescription: "The unique identifier of the export image task.",
							Computed:            true,
						},
						"image_arn": schema.StringAttribute{
							Description:         "ARN of the exported AppStream image.",
							MarkdownDescription: "The Amazon Resource Name (ARN) of the AppStream Applications image being exported.",
							Computed:            true,
						},
						"ami_name": schema.StringAttribute{
							Description:         "AMI name.",
							MarkdownDescription: "The name of the EC2 AMI created by the export image task.",
							Computed:            true,
						},
						"ami_description": schema.StringAttribute{
							Description:         "AMI description.",
							MarkdownDescription: "The description applied to the exported EC2 AMI, if set.",
							Computed:            true,
						},
						"ami_id": schema.StringAttribute{
							Description: "AMI ID.",
							MarkdownDescription: "The ID of the EC2 AMI created by the export image task. " +
								"This attribute is only populated when the task completes successfully.",
							Computed: true,
						},
						"state": schema.StringAttribute{
							Description: "Export task state.",
							MarkdownDescription: "The current state of the export image task. " +
								"Valid values are `EXPORTING`, `COMPLETED`, or `FAILED`.",
							Computed: true,
						},
						"created_date": schema.StringAttribute{
							Description:         "Task creation time.",
							MarkdownDescription: "The timestamp when the export image task was created, in RFC 3339 format.",
							Computed:            true,
						},
						"tag_specifications": schema.MapAttribute{
							Description:         "AMI tag specifications.",
							MarkdownDescription: "The tags applied to the exported EC2 AMI.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"error_details": schema.SetNestedAttribute{
							Description: "Export task error details.",
							MarkdownDescription: "Details about errors that occurred during the export process. " +
								"This attribute is only populated when the task state is `FAILED`.",
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"error_code": schema.StringAttribute{
										Description:         "Error code.",
										MarkdownDescription: "The error code reported by AWS.",
										Computed:            true,
									},
									"error_message": schema.StringAttribute{
										Description:         "Error message.",
										MarkdownDescription: "The human-readable error message reported by AWS.",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
