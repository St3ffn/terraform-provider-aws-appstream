// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription

import "github.com/hashicorp/terraform-plugin-framework/types"

type dataSourceModel struct {
	// S3BucketName is the Amazon S3 bucket where generated usage reports are stored (computed).
	S3BucketName types.String `tfsdk:"s3_bucket_name"`
	// Schedule is the schedule for generating usage reports (computed).
	// Currently, only DAILY is supported.
	Schedule types.String `tfsdk:"schedule"`
	// LastGeneratedReportDate is the timestamp when the last usage report was generated (computed).
	LastGeneratedReportDate types.String `tfsdk:"last_generated_report_date"`
	// SubscriptionErrors is the list of errors reported by AWS for the usage report subscription (computed).
	SubscriptionErrors types.Set `tfsdk:"subscription_errors"`
}
