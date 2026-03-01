// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/st3ffn/terraform-provider-aws-appstream/internal/util"
)

var subscriptionErrorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"error_code":    types.StringType,
		"error_message": types.StringType,
	},
}

func flattenSubscriptionErrors(
	ctx context.Context, awsErrors []awstypes.LastReportGenerationExecutionError, diags *diag.Diagnostics,
) types.Set {

	if len(awsErrors) == 0 {
		return types.SetNull(subscriptionErrorObjectType)
	}

	out := make([]resourceModelSubscriptionErrors, 0, len(awsErrors))
	for _, e := range awsErrors {
		out = append(out, resourceModelSubscriptionErrors{
			ErrorCode:    types.StringValue(string(e.ErrorCode)),
			ErrorMessage: util.StringOrNull(e.ErrorMessage),
		})
	}

	setVal, d := types.SetValueFrom(ctx, subscriptionErrorObjectType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(subscriptionErrorObjectType)
	}

	return setVal
}
