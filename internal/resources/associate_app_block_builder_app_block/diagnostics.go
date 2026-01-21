// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_app_block_builder_app_block

import "github.com/hashicorp/terraform-plugin-framework/diag"

type diagnosticMode string

const (
	diagnosticPlan   diagnosticMode = "plan"
	diagnosticRead   diagnosticMode = "read"
	diagnosticDelete diagnosticMode = "delete"
)

func addDiagnostics(model model, diags *diag.Diagnostics, mode diagnosticMode) {
	if model.AppBlockBuilderName.IsNull() || model.AppBlockBuilderName.IsUnknown() ||
		model.AppBlockARN.IsNull() || model.AppBlockARN.IsUnknown() {

		switch mode {
		case diagnosticPlan:
			diags.AddError(
				"Invalid Terraform Plan",
				"Cannot associate app block builder with app block because "+
					"app_block_builder_name and app_block_arn must be known.",
			)

		case diagnosticDelete:
			diags.AddError(
				"Invalid Terraform State",
				"Cannot disassociate app block builder from app block because "+
					"app_block_builder_name and app_block_arn must be known.",
			)

		case diagnosticRead:
			diags.AddError(
				"Invalid Terraform State",
				"Required attributes app_block_builder_name and app_block_arn are missing from state. "+
					"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
			)
		}
	}
}
