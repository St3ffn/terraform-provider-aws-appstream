// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Schema returns the static Terraform schema for this component.
// It defines supported attributes/blocks plus validators and plan behavior.
func (ds *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read AWS AppStream Usage Report Subscription",
		MarkdownDescription: "Reads the AppStream usage report subscription for the current AWS account and region. " +
			"This data source is read-only and can be used to inspect whether usage reporting is enabled, " +
			"the S3 bucket where reports are delivered, and any errors reported by AWS. " +
			"Only one usage report subscription can exist per AWS account and region.",
		Attributes: map[string]schema.Attribute{
			"s3_bucket_name": schema.StringAttribute{
				Description: "S3 bucket where usage reports are stored.",
				MarkdownDescription: "The Amazon S3 bucket where generated usage reports are stored. " +
					"If on-instance session scripts and S3 logging are enabled, the same bucket is reused.",
				Computed: true,
			},
			"schedule": schema.StringAttribute{
				Description:         "Schedule for generating usage reports.",
				MarkdownDescription: "The schedule for generating usage reports. Currently, only `DAILY` is supported.",
				Computed:            true,
			},
			"last_generated_report_date": schema.StringAttribute{
				Description:         "Timestamp of the last generated usage report.",
				MarkdownDescription: "The time when the last usage report was generated, in RFC 3339 format.",
				Computed:            true,
			},
			"subscription_errors": schema.SetNestedAttribute{
				Description:         "Errors reported for the Usage Report Subscription.",
				MarkdownDescription: "Errors returned by AWS if usage reports could not be generated.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							Description:         "Subscription error code.",
							MarkdownDescription: "The error code returned by AWS.",
							Computed:            true,
						},
						"error_message": schema.StringAttribute{
							Description:         "Subscription error message.",
							MarkdownDescription: "The human-readable error message returned by AWS.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
