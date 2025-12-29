// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var attributeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	},
}

func flattenAttributes(
	ctx context.Context, awsEntitlementAttributes []awstypes.EntitlementAttribute, diags *diag.Diagnostics,
) types.Set {

	attrs := make([]attributeModel, 0, len(awsEntitlementAttributes))
	for _, a := range awsEntitlementAttributes {
		attrs = append(attrs, attributeModel{
			Name:  types.StringValue(aws.ToString(a.Name)),
			Value: types.StringValue(aws.ToString(a.Value)),
		})
	}

	setVal, d := types.SetValueFrom(ctx, attributeObjectType, attrs)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(attributeObjectType)
	}

	return setVal
}
