// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type entitlementModel struct {
	ID            types.String `tfsdk:"id"`
	StackName     types.String `tfsdk:"stack_name"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	AppVisibility types.String `tfsdk:"app_visibility"`
	Attributes    types.Set    `tfsdk:"attributes"`
}

type entitlementAttributeModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
