// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Read(ctx context.Context, req tfresource.ReadRequest, resp *tfresource.ReadResponse) {
	var state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ctx.Err(); err != nil {
		return
	}

	newState, diags := r.readUsageSubscription(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		if ctx.Err() != nil {
			return
		}
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *resource) readUsageSubscription(ctx context.Context) (*resourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	out, err := r.appstreamClient.DescribeUsageReportSubscriptions(ctx, &awsappstream.DescribeUsageReportSubscriptionsInput{})
	if err != nil {
		if util.IsContextCanceled(err) {
			return nil, diags
		}

		if util.IsAppStreamNotFound(err) {
			return nil, diags
		}
		diags.AddError(
			"Error Reading AWS AppStream Usage Report Subscription",
			fmt.Sprintf("Could not read usage report subscription: %v", err),
		)
		return nil, diags
	}

	if len(out.UsageReportSubscriptions) == 0 {
		return nil, diags
	}

	subscription := out.UsageReportSubscriptions[0]
	state := &resourceModel{
		ID:                      types.StringValue(usageReportSubscriptionID),
		S3BucketName:            util.StringOrNull(subscription.S3BucketName),
		Schedule:                types.StringValue(string(subscription.Schedule)),
		LastGeneratedReportDate: util.StringFromTime(subscription.LastGeneratedReportDate),
		SubscriptionErrors:      flattenSubscriptionErrors(ctx, subscription.SubscriptionErrors, &diags),
	}

	if diags.HasError() {
		return nil, diags
	}
	return state, diags
}
