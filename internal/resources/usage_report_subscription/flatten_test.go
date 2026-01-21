// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package usage_report_subscription

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenSubscriptionErrors(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		in   []awstypes.LastReportGenerationExecutionError
		want types.Set
	}{
		{
			name: "empty_slice",
			in:   nil,
			want: types.SetNull(subscriptionErrorObjectType),
		},
		{
			name: "single_error",
			in: []awstypes.LastReportGenerationExecutionError{
				{
					ErrorCode:    awstypes.UsageReportExecutionErrorCodeAccessDenied,
					ErrorMessage: aws.String("boom"),
				},
			},
			want: types.SetValueMust(
				subscriptionErrorObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						subscriptionErrorObjectType.AttrTypes,
						map[string]attr.Value{
							"error_code":    types.StringValue(string(awstypes.UsageReportExecutionErrorCodeAccessDenied)),
							"error_message": types.StringValue("boom"),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := flattenSubscriptionErrors(ctx, tt.in, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
