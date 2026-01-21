// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription

import "github.com/hashicorp/terraform-plugin-framework/types"

type resourceModel struct {
	// ID is a synthetic identifier for the usage report subscription.
	ID types.String `tfsdk:"id"`
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

type subscriptionErrorModel struct {
	// ErrorCode is the error code returned by AWS when a usage report could not be generated (computed).
	ErrorCode types.String `tfsdk:"error_code"`
	// ErrorMessage is the human-readable error message returned by AWS (computed).
	ErrorMessage types.String `tfsdk:"error_message"`
}
