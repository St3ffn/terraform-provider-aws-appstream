// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/diag"

type assocDiagMode string

const (
	assocDiagPlan   assocDiagMode = "plan"
	assocDiagRead   assocDiagMode = "read"
	assocDiagDelete assocDiagMode = "delete"
)

func addAssocPartsDiagnostics(m associateApplicationEntitlementModel, diags *diag.Diagnostics, mode assocDiagMode) {
	if m.StackName.IsNull() || m.StackName.IsUnknown() ||
		m.EntitlementName.IsNull() || m.EntitlementName.IsUnknown() ||
		m.ApplicationIdentifier.IsNull() || m.ApplicationIdentifier.IsUnknown() {

		switch mode {
		case assocDiagPlan:
			diags.AddError(
				"Invalid Terraform Plan",
				"Cannot associate application to entitlement because stack_name, entitlement_name, and application_identifier must be known.",
			)
		case assocDiagDelete:
			diags.AddError(
				"Invalid Terraform State",
				"Cannot disassociate application from entitlement because stack_name, entitlement_name, and application_identifier must be known.",
			)
		case assocDiagRead:
			diags.AddError(
				"Invalid Terraform State",
				"Required attributes stack_name, entitlement_name, and application_identifier are missing from state. "+
					"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
			)
		}
	}
}
