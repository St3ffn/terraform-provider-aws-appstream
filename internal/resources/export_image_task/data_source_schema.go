// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package export_image_task

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (ds *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read an AWS AppStream Export Image Task",
		MarkdownDescription: "Reads a single AppStream export image task. " +
			"An export image task represents a long-running operation that exports a AppStream Applications image " +
			"to an EC2 AMI. This data source can be used to monitor task progress, retrieve the resulting AMI ID, " +
			"or inspect error details if the export fails.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the export image task.",
				MarkdownDescription: "A synthetic identifier for the export image task, equal to the task ID.",
				Computed:            true,
			},
			"task_id": schema.StringAttribute{
				Description:         "Export image task ID.",
				MarkdownDescription: "The unique identifier of the export image task.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`),
						"must match [a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}",
					),
				},
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
	}
}
