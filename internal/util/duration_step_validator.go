// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type durationStepValidator struct {
	min       int32
	max       int32
	step      int32
	allowZero bool
}

// Description returns a human-readable description of the validator behavior.
func (v durationStepValidator) Description(_ context.Context) string {
	if v.allowZero {
		return fmt.Sprintf(
			"value must be 0 or between %d and %d in increments of %d seconds", v.min, v.max, v.step,
		)
	}

	return fmt.Sprintf("value must be between %d and %d in increments of %d seconds", v.min, v.max, v.step)
}

// MarkdownDescription returns a Markdown validator description for Terraform docs.
func (v durationStepValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateInt32 validates int32 duration values against configured bounds and step.
func (v durationStepValidator) ValidateInt32(
	ctx context.Context, req validator.Int32Request, resp *validator.Int32Response,
) {

	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueInt32()

	if v.allowZero && value == 0 {
		return
	}

	if value < v.min || value > v.max {
		resp.Diagnostics.Append(
			validatordiag.InvalidAttributeValueMatchDiagnostic(
				req.Path,
				v.Description(ctx),
				fmt.Sprintf("%d", value),
			),
		)
		return
	}

	if v.step > 0 && value%v.step != 0 {
		resp.Diagnostics.Append(
			validatordiag.InvalidAttributeValueMatchDiagnostic(
				req.Path,
				v.Description(ctx),
				fmt.Sprintf("%d", value),
			),
		)
	}
}

// DurationWithStep returns an int32 validator for duration values with step constraints.
func DurationWithStep(minValue, maxValue, step int32, allowZero bool) validator.Int32 {
	return durationStepValidator{
		min:       minValue,
		max:       maxValue,
		step:      step,
		allowZero: allowZero,
	}
}
