// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"errors"

	"github.com/aws/smithy-go"
)

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

func IsOperationNotPermittedException(err error) bool {
	return IsAWSAPIError(err, "OperationNotPermittedException")
}

func IsResourceNotFoundException(err error) bool {
	return IsAWSAPIError(err, "ResourceNotFoundException")
}

func IsConcurrentModificationException(err error) bool {
	return IsAWSAPIError(err, "ConcurrentModificationException")
}

func IsEntitlementNotFoundException(err error) bool {
	return IsAWSAPIError(err, "EntitlementNotFoundException")
}

func IsResourceNotAvailableException(err error) bool {
	return IsAWSAPIError(err, "ResourceNotAvailableException")
}

func IsResourceAlreadyExists(err error) bool {
	return IsAWSAPIError(err, "ResourceAlreadyExistsException")
}

func IsResourceInUseException(err error) bool {
	return IsAWSAPIError(err, "ResourceInUseException")
}

func IsEntitlementAlreadyExists(err error) bool {
	return IsAWSAPIError(err, "EntitlementAlreadyExistsException")
}

func IsAppStreamNotFound(err error) bool {
	return IsAWSAPIError(err, "ResourceNotFoundException", "EntitlementNotFoundException")
}

func AWSEnumToSlice[T ~string](awsEnumValuesFunc func(T) []T) []string {
	var zero T
	enumValues := awsEnumValuesFunc(zero)

	result := make([]string, 0, len(enumValues))
	for _, value := range enumValues {
		result = append(result, string(value))
	}

	return result
}
