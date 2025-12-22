// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/diag"

func addAssociateApplicationEntitlementDiagnostics(
	model associateApplicationEntitlementModel, diags *diag.Diagnostics, mode associateDiagnosticMode,
) {

	if model.StackName.IsNull() || model.StackName.IsUnknown() ||
		model.EntitlementName.IsNull() || model.EntitlementName.IsUnknown() ||
		model.ApplicationIdentifier.IsNull() || model.ApplicationIdentifier.IsUnknown() {

		switch mode {
		case associateDiagnosticPlan:
			diags.AddError(
				"Invalid Terraform Plan",
				"Cannot associate application to entitlement because stack_name, entitlement_name, and application_identifier must be known.",
			)
		case associateDiagnosticDelete:
			diags.AddError(
				"Invalid Terraform State",
				"Cannot disassociate application from entitlement because stack_name, entitlement_name, and application_identifier must be known.",
			)
		case associateDiagnosticRead:
			diags.AddError(
				"Invalid Terraform State",
				"Required attributes stack_name, entitlement_name, and application_identifier are missing from state. "+
					"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
			)
		}
	}
}
