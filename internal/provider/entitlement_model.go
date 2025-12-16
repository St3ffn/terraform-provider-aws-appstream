// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type entitlementModel struct {
	ID               types.String `tfsdk:"id"`
	StackName        types.String `tfsdk:"stack_name"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	AppVisibility    types.String `tfsdk:"app_visibility"`
	Attributes       types.Set    `tfsdk:"attributes"`
	CreatedTime      types.String `tfsdk:"created_time"`
	LastModifiedTime types.String `tfsdk:"last_modified_time"`
}

type entitlementAttributeModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
