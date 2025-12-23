// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func expandEntitlementAttributes(
	ctx context.Context, tfAttributes types.Set, diags *diag.Diagnostics,
) []awstypes.EntitlementAttribute {

	var attrs []entitlementAttributeModel
	diags.Append(tfAttributes.ElementsAs(ctx, &attrs, false)...)
	if diags.HasError() {
		return nil
	}

	awsAttrs := make([]awstypes.EntitlementAttribute, 0, len(attrs))
	for _, a := range attrs {
		if a.Name.IsNull() || a.Name.IsUnknown() ||
			a.Value.IsNull() || a.Value.IsUnknown() {

			diags.AddError(
				"Invalid Terraform Plan",
				"All entitlement attributes must have known, non-null name and value.",
			)
			return nil
		}

		awsAttrs = append(awsAttrs, awstypes.EntitlementAttribute{
			Name:  aws.String(a.Name.ValueString()),
			Value: aws.String(a.Value.ValueString()),
		})
	}

	return awsAttrs
}
