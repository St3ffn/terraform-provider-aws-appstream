// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ValidateRegexPattern validates that pattern compiles as a regular expression.
func ValidateRegexPattern(pattern string) error {
	_, err := regexp.CompilePOSIX(pattern)

	if err != nil {
		return fmt.Errorf("invalid regular expression: %s", err.Error())
	}

	return nil
}

type regexValidator struct{}

// Description returns a human-readable description of the validator behavior.
func (v regexValidator) Description(_ context.Context) string {
	return "string must be a valid regular expression"
}

// MarkdownDescription returns a Markdown validator description for Terraform docs.
func (v regexValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString validates a Terraform string value for the current validator.
func (v regexValidator) ValidateString(
	ctx context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {

	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if err := ValidateRegexPattern(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.Append(
			validatordiag.InvalidAttributeValueMatchDiagnostic(
				req.Path,
				v.Description(ctx),
				req.ConfigValue.ValueString(),
			),
		)
	}
}

// ValidRegex returns a string validator that accepts valid regular expression patterns.
func ValidRegex() validator.String {
	return regexValidator{}
}
