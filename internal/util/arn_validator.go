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

func ValidateARNValue(value, service, resourcePrefix string) error {
	parsed, err := awsarn.Parse(value)
	if err != nil {
		return fmt.Errorf("invalid ARN")
	}

	if service != "" && parsed.Service != service {
		return fmt.Errorf("expected ARN service %q", service)
	}

	if resourcePrefix != "" && !strings.HasPrefix(parsed.Resource, resourcePrefix) {
		return fmt.Errorf("expected ARN resource prefix %q", resourcePrefix)
	}

	return nil
}

type arnValidator struct {
	service        string
	resourcePrefix string
}

func (v arnValidator) Description(_ context.Context) string {
	if v.service == "" && v.resourcePrefix == "" {
		return "string must be a valid ARN"
	}
	if v.resourcePrefix == "" {
		return fmt.Sprintf("string must be a valid ARN for service %s", v.service)
	}
	return fmt.Sprintf("string must be a valid ARN for service %s with resource prefix %s",
		v.service, v.resourcePrefix)
}

func (v arnValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v arnValidator) ValidateString(
	ctx context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {

	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if err := ValidateARNValue(req.ConfigValue.ValueString(), v.service, v.resourcePrefix); err != nil {
		resp.Diagnostics.Append(
			validatordiag.InvalidAttributeValueMatchDiagnostic(
				req.Path,
				v.Description(ctx),
				req.ConfigValue.ValueString(),
			),
		)
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
