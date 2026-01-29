// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription

import (
	"context"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an AWS AppStream Usage Report Subscription.",
		MarkdownDescription: "Manages an AppStream usage report subscription. " +
			"Usage reports are generated daily and delivered to an Amazon S3 bucket. " +
			"Only one usage report subscription can exist per AWS account and region.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier of the usage report subscription.",
				MarkdownDescription: "A synthetic identifier for the usage report subscription.",
				Computed:            true,
			},
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
				Description:         "Errors reported for the usage report subscription.",
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
