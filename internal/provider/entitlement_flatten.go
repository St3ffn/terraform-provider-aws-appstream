// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var entitlementAttributeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	},
}

func flattenEntitlementAttributes(
	ctx context.Context, awsAttrs []awstypes.EntitlementAttribute, diags *diag.Diagnostics,
) types.Set {

	attrs := make([]entitlementAttributeModel, 0, len(awsAttrs))
	for _, a := range awsAttrs {
		attrs = append(attrs, entitlementAttributeModel{
			Name:  types.StringValue(aws.ToString(a.Name)),
			Value: types.StringValue(aws.ToString(a.Value)),
		})
	}

	setVal, d := types.SetValueFrom(ctx, entitlementAttributeObjectType, attrs)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(entitlementAttributeObjectType)
	}

	return setVal
}
