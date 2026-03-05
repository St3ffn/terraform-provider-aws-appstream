// Copyright St3ffn 2025, 2026
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

// ValidateARNValue validates an ARN against expected service and resource prefix constraints.
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

// Description returns a human-readable description of the validator behavior.
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

// MarkdownDescription returns a Markdown validator description for Terraform docs.
func (v arnValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString validates a Terraform string value for the current validator.
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

// ValidARN returns a string validator for generic ARN values.
func ValidARN() validator.String {
	return arnValidator{}
}

// ValidARNWithService returns an ARN validator constrained to a specific AWS service.
func ValidARNWithService(service string) validator.String {
	return arnValidator{service: service}
}

// ValidARNWithServiceAndResource returns an ARN validator constrained to service and resource prefix.
func ValidARNWithServiceAndResource(service, resourcePrefix string) validator.String {
	return arnValidator{
		service:        service,
		resourcePrefix: resourcePrefix,
	}
}
