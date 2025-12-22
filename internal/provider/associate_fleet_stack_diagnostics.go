// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/diag"

func addAssociateFleetStackDiagnostics(
	model associateFleetStackModel, diags *diag.Diagnostics, mode associateDiagnosticMode,
) {
	if model.FleetName.IsNull() || model.FleetName.IsUnknown() ||
		model.StackName.IsNull() || model.StackName.IsUnknown() {

		switch mode {
		case associateDiagnosticPlan:
			diags.AddError(
				"Invalid Terraform Plan",
				"Cannot associate fleet to stack because fleet_name and stack_name must be known.",
			)
		case associateDiagnosticDelete:
			diags.AddError(
				"Invalid Terraform State",
				"Cannot disassociate fleet from stack because fleet_name and stack_name must be known.",
			)
		case associateDiagnosticRead:
			diags.AddError(
				"Invalid Terraform State",
				"Required attributes fleet_name and stack_name are missing from state. "+
					"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
			)
		}
	}
}
