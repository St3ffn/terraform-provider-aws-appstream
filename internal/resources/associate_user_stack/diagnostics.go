// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack

import "github.com/hashicorp/terraform-plugin-framework/diag"

type diagnosticMode string

const (
	diagnosticPlan   diagnosticMode = "plan"
	diagnosticRead   diagnosticMode = "read"
	diagnosticDelete diagnosticMode = "delete"
)

func addDiagnostics(model model, diags *diag.Diagnostics, mode diagnosticMode) {
	if model.StackName.IsNull() || model.StackName.IsUnknown() ||
		model.UserName.IsNull() || model.UserName.IsUnknown() ||
		model.AuthenticationType.IsNull() || model.AuthenticationType.IsUnknown() {

		switch mode {
		case diagnosticPlan:
			diags.AddError(
				"Invalid Terraform Plan",
				"Cannot associate user to stack because stack_name, user_name, and authentication_type must be known.",
			)

		case diagnosticDelete:
			diags.AddError(
				"Invalid Terraform State",
				"Cannot disassociate user from stack because stack_name, user_name, and authentication_type must be known.",
			)

		case diagnosticRead:
			diags.AddError(
				"Invalid Terraform State",
				"Required attributes stack_name, user_name, and authentication_type are missing from state. "+
					"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
			)
		}
	}
}
