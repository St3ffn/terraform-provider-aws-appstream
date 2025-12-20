// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type associateApplicationEntitlementModel struct {
	// ID is a synthetic identifier composed of "<stack_name>|<entitlement_name>|<application_identifier>"
	ID types.String `tfsdk:"id"`
	// StackName is the name of the AppStream stack that owns the entitlement (required).
	StackName types.String `tfsdk:"stack_name"`
	// EntitlementName is the name of the entitlement to which the application is associated (required).
	EntitlementName types.String `tfsdk:"entitlement_name"`
	// ApplicationIdentifier is the identifier of the AppStream application being associated (required).
	ApplicationIdentifier types.String `tfsdk:"application_identifier"`
}
