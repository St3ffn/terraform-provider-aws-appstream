// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ValidateRegexPattern(pattern string) error {
	_, err := regexp.CompilePOSIX(pattern)

	if err != nil {
		return fmt.Errorf("invalid regular expression: %s", err.Error())
	}

	return nil
}

type regexValidator struct{}

func (v regexValidator) Description(_ context.Context) string {
	return "string must be a valid regular expression"
}

func (v regexValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

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

func ValidRegex() validator.String {
	return regexValidator{}
}
