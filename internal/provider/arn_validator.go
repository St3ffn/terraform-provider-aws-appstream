// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type arnValidator struct{}

func (a arnValidator) Description(_ context.Context) string {
	return "string must be a valid ARN"
}

func (a arnValidator) MarkdownDescription(_ context.Context) string {
	return "string must be a valid ARN"
}

func (a arnValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	if _, err := arn.Parse(value); err != nil {
		resp.Diagnostics.Append(
			validatordiag.InvalidAttributeValueMatchDiagnostic(
				req.Path,
				a.Description(ctx),
				req.ConfigValue.ValueString(),
			),
		)
	}
}

func validARN() arnValidator {
	return arnValidator{}
}
