// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type entitlementModel struct {
	// ID is a synthetic identifier composed of "<stack_name>|<entitlement_name>".
	ID types.String `tfsdk:"id"`
	// StackName is the name of the AppStream stack in which the entitlement is defined (required).
	StackName types.String `tfsdk:"stack_name"`
	// Name is the name of the AppStream entitlement within the stack (required).
	Name types.String `tfsdk:"name"`
	// Description is the description of the entitlement (optional).
	Description types.String `tfsdk:"description"`
	// AppVisibility controls which applications are visible to users matching this entitlement (required).
	AppVisibility types.String `tfsdk:"app_visibility"`
	// Attributes is the set of attribute rules used to match federated user attributes (required).
	Attributes types.Set `tfsdk:"attributes"`
	// CreatedTime is the timestamp when the entitlement was created (computed).
	CreatedTime types.String `tfsdk:"created_time"`
	// LastModifiedTime is the timestamp when the entitlement was last modified (computed).
	LastModifiedTime types.String `tfsdk:"last_modified_time"`
}

type entitlementAttributeModel struct {
	// Name is the name of the entitlement attribute (required).
	Name types.String `tfsdk:"name"`
	// Value is the value of the entitlement attribute (required).
	Value types.String `tfsdk:"value"`
}
