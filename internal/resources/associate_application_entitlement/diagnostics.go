// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_application_entitlement

import "github.com/hashicorp/terraform-plugin-framework/diag"

type diagnosticMode string

const (
	diagnosticPlan   diagnosticMode = "plan"
	diagnosticRead   diagnosticMode = "read"
	diagnosticDelete diagnosticMode = "delete"
)

func addDiagnostics(model model, diags *diag.Diagnostics, mode diagnosticMode) {
	if model.StackName.IsNull() || model.StackName.IsUnknown() ||
		model.EntitlementName.IsNull() || model.EntitlementName.IsUnknown() ||
		model.ApplicationIdentifier.IsNull() || model.ApplicationIdentifier.IsUnknown() {

		switch mode {
		case diagnosticPlan:
			diags.AddError(
				"Invalid Terraform Plan",
				"Cannot associate application to entitlement because stack_name, entitlement_name, and application_identifier must be known.",
			)
		case diagnosticDelete:
			diags.AddError(
				"Invalid Terraform State",
				"Cannot disassociate application from entitlement because stack_name, entitlement_name, and application_identifier must be known.",
			)
		case diagnosticRead:
			diags.AddError(
				"Invalid Terraform State",
				"Required attributes stack_name, entitlement_name, and application_identifier are missing from state. "+
					"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
			)
		}
	}
}
