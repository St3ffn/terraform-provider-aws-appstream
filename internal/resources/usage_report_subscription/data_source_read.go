// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

// Read fetches the current remote object and refreshes Terraform state from AWS.
// When the object no longer exists remotely, the resource is removed from state
// to converge Terraform with external deletions.
func (ds *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config dataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	out, err := ds.appstreamClient.DescribeUsageReportSubscriptions(ctx, &awsappstream.DescribeUsageReportSubscriptionsInput{})
	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		if util.IsAppStreamNotFound(err) {
			resp.Diagnostics.AddError(
				"AWS AppStream Usage Report Subscription Not Found",
				"No usage report subscription was found.",
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AWS AppStream Usage Report Subscription",
			fmt.Sprintf("Could not read usage report subscription: %v", err),
		)
		return
	}

	if len(out.UsageReportSubscriptions) == 0 {
		resp.Diagnostics.AddError(
			"AWS AppStream Usage Report Subscription Not Found",
			"No usage report subscription was found.",
		)
		return
	}

	subscription := out.UsageReportSubscriptions[0]

	state := &dataSourceModel{
		S3BucketName:            util.StringOrNull(subscription.S3BucketName),
		Schedule:                types.StringValue(string(subscription.Schedule)),
		LastGeneratedReportDate: util.StringFromTime(subscription.LastGeneratedReportDate),
		SubscriptionErrors:      flattenSubscriptionErrors(ctx, subscription.SubscriptionErrors, &resp.Diagnostics),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
