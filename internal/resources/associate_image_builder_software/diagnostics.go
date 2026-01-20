// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

import "github.com/hashicorp/terraform-plugin-framework/diag"

type diagnosticMode string

const (
	diagnosticPlan   diagnosticMode = "plan"
	diagnosticRead   diagnosticMode = "read"
	diagnosticDelete diagnosticMode = "delete"
	diagnosticUpdate diagnosticMode = "update"
)

func addDiagnostics(model model, diags *diag.Diagnostics, mode diagnosticMode) {
	if model.ImageBuilderARN.IsNull() || model.ImageBuilderARN.IsUnknown() {

		switch mode {
		case diagnosticPlan:
			diags.AddError(
				"Invalid Terraform Plan",
				"Cannot associate software to image builder because image_builder_arn must be known.",
			)
		case diagnosticDelete:
			diags.AddError(
				"Invalid Terraform State",
				"Cannot disassociate software from image builder because image_builder_arn must be known.",
			)
		case diagnosticRead:
			diags.AddError(
				"Invalid Terraform State",
				"Required attribute image_builder_arn is missing from state. "+
					"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
			)
		case diagnosticUpdate:
			diags.AddError(
				"Invalid Terraform Plan",
				"Cannot update association because image_builder_arn must be known.",
			)
		}
	}
}
