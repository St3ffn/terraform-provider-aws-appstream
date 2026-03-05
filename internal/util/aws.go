// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"errors"

	"github.com/aws/smithy-go"
)

// IsAWSAPIError reports whether err matches one of the given AWS API error codes.
func IsAWSAPIError(err error, code ...string) bool {
	var apiErr smithy.APIError
	if err == nil || !errors.As(err, &apiErr) {
		return false
	}

	for _, c := range code {
		if apiErr.ErrorCode() == c {
			return true
		}
	}
	return false
}

// IsOperationNotPermittedException reports whether err is an OperationNotPermittedException.
func IsOperationNotPermittedException(err error) bool {
	return IsAWSAPIError(err, "OperationNotPermittedException")
}

// IsResourceNotFoundException reports whether err is a ResourceNotFoundException.
func IsResourceNotFoundException(err error) bool {
	return IsAWSAPIError(err, "ResourceNotFoundException")
}

// IsConcurrentModificationException reports whether err is a ConcurrentModificationException.
func IsConcurrentModificationException(err error) bool {
	return IsAWSAPIError(err, "ConcurrentModificationException")
}

// IsEntitlementNotFoundException reports whether err is an EntitlementNotFoundException.
func IsEntitlementNotFoundException(err error) bool {
	return IsAWSAPIError(err, "EntitlementNotFoundException")
}

// IsResourceNotAvailableException reports whether err is a ResourceNotAvailableException.
func IsResourceNotAvailableException(err error) bool {
	return IsAWSAPIError(err, "ResourceNotAvailableException")
}

// IsInvalidRoleException reports whether err is an InvalidRoleException.
func IsInvalidRoleException(err error) bool {
	return IsAWSAPIError(err, "InvalidRoleException")
}

// IsResourceAlreadyExists reports whether err is a ResourceAlreadyExistsException.
func IsResourceAlreadyExists(err error) bool {
	return IsAWSAPIError(err, "ResourceAlreadyExistsException")
}

// IsResourceInUseException reports whether err is a ResourceInUseException.
func IsResourceInUseException(err error) bool {
	return IsAWSAPIError(err, "ResourceInUseException")
}

// IsEntitlementAlreadyExists reports whether err is an EntitlementAlreadyExistsException.
func IsEntitlementAlreadyExists(err error) bool {
	return IsAWSAPIError(err, "EntitlementAlreadyExistsException")
}

// IsAppStreamNotFound reports whether err indicates a missing AppStream resource.
func IsAppStreamNotFound(err error) bool {
	return IsAWSAPIError(err, "ResourceNotFoundException", "EntitlementNotFoundException")
}

// AWSEnumToSlice converts an AWS SDK enum value list function into a []string.
func AWSEnumToSlice[T ~string](awsEnumValuesFunc func(T) []T) []string {
	var zero T
	enumValues := awsEnumValuesFunc(zero)

	result := make([]string, 0, len(enumValues))
	for _, value := range enumValues {
		result = append(result, string(value))
	}

	return result
}
