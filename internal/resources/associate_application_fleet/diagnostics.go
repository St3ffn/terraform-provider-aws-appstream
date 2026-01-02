// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_application_fleet

import "github.com/hashicorp/terraform-plugin-framework/diag"

type diagnosticMode string

const (
	diagnosticPlan   diagnosticMode = "plan"
	diagnosticRead   diagnosticMode = "read"
	diagnosticDelete diagnosticMode = "delete"
)

func addDiagnostics(model model, diags *diag.Diagnostics, mode diagnosticMode) {
	if model.FleetName.IsNull() || model.FleetName.IsUnknown() ||
		model.ApplicationARN.IsNull() || model.ApplicationARN.IsUnknown() {

		switch mode {
		case diagnosticPlan:
			diags.AddError(
				"Invalid Terraform Plan",
				"Cannot associate application to fleet because fleet_name and application_arn must be known.",
			)
		case diagnosticDelete:
			diags.AddError(
				"Invalid Terraform State",
				"Cannot disassociate application from fleet because fleet_name and application_arn must be known.",
			)
		case diagnosticRead:
			diags.AddError(
				"Invalid Terraform State",
				"Required attributes fleet_name and application_arn are missing from state. "+
					"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
			)
		}
	}
}
