// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription

import (
	"context"
	"fmt"

	awsappstream "github.com/aws/aws-sdk-go-v2/service/appstream"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

func (r *resource) Create(ctx context.Context, req tfresource.CreateRequest, resp *tfresource.CreateResponse) {
	var plan resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := ctx.Err(); err != nil {
		return
	}

	_, err := r.appstreamClient.CreateUsageReportSubscription(ctx, &awsappstream.CreateUsageReportSubscriptionInput{})

	if err != nil {
		if util.IsContextCanceled(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error Creating AWS AppStream Usage Report Subscription",
			fmt.Sprintf("Could not create usage report subscription: %v", err),
		)
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
