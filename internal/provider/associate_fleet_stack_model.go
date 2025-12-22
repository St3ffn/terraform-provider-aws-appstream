// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type associateFleetStackModel struct {
	// ID is a synthetic identifier composed of "<fleet_name>|<stack_name>"
	ID types.String `tfsdk:"id"`
	// FleetName is the name of the AppStream fleet to be associated with the stack (required)
	FleetName types.String `tfsdk:"fleet_name"`
	// StackName is the name of the AppStream stack to associate with the fleet (required).
	StackName types.String `tfsdk:"stack_name"`
}
