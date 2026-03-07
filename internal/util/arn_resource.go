// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsarn "github.com/aws/aws-sdk-go-v2/aws/arn"
)

// ARNResourceSuffixOrNil parses an ARN string pointer and returns the resource
// suffix after resourcePrefix when service and resource prefix match. It returns
// nil for nil input, parse failures, mismatches, or empty suffix values.
func ARNResourceSuffixOrNil(value *string, service, resourcePrefix string) *string {
	// expected arn:aws:service:<region>:<account>:resource_prefix/<name>
	if value == nil {
		return nil
	}

	parsed, err := awsarn.Parse(*value)
	if err != nil {
		return nil
	}

	if parsed.Service != service || !strings.HasPrefix(parsed.Resource, resourcePrefix) {
		return nil
	}

	name := strings.TrimPrefix(parsed.Resource, resourcePrefix)
	if name == "" {
		return nil
	}

	return aws.String(name)
}
