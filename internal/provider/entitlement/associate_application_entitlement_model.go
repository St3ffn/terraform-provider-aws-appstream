// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type associateApplicationEntitlementModel struct {
	ID                    types.String `tfsdk:"id"`
	StackName             types.String `tfsdk:"stack_name"`
	EntitlementName       types.String `tfsdk:"entitlement_name"`
	ApplicationIdentifier types.String `tfsdk:"application_identifier"`
}
