// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"fmt"
	"strings"

	awsarn "github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type arnValidator struct {
	service        string
	resourcePrefix string
}

func (a arnValidator) Description(_ context.Context) string {
	if a.service == "" && a.resourcePrefix == "" {
		return "string must be a valid ARN"
	}
	if a.resourcePrefix == "" {
		return fmt.Sprintf("string must be a valid ARN for service %s", a.service)
	}
	return fmt.Sprintf("string must be a valid ARN for service %s with resource prefix %s",
		a.service, a.resourcePrefix)
}

func (a arnValidator) MarkdownDescription(ctx context.Context) string {
	return a.Description(ctx)
}

func (a arnValidator) ValidateString(
	ctx context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {

	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	parsed, err := awsarn.Parse(value)
	if err != nil {
		resp.Diagnostics.Append(
			validatordiag.InvalidAttributeValueMatchDiagnostic(
				req.Path,
				a.Description(ctx),
				req.ConfigValue.ValueString(),
			),
		)
		return
	}

	if a.service != "" && parsed.Service != a.service {
		resp.Diagnostics.Append(
			validatordiag.InvalidAttributeValueMatchDiagnostic(
				req.Path,
				a.Description(ctx),
				value,
			),
		)
		return
	}

	if a.resourcePrefix != "" && !strings.HasPrefix(parsed.Resource, a.resourcePrefix) {
		resp.Diagnostics.Append(
			validatordiag.InvalidAttributeValueMatchDiagnostic(
				req.Path,
				a.Description(ctx),
				value,
			),
		)
		return
	}
}

//nolint:unused // kept for future validators
func ValidARN() validator.String {
	return arnValidator{}
}

//nolint:unused // kept for future validators
func ValidARNWithService(service string) validator.String {
	return arnValidator{service: service}
}

func ValidARNWithServiceAndResource(service, resourcePrefix string) validator.String {
	return arnValidator{
		service:        service,
		resourcePrefix: resourcePrefix,
	}
}
