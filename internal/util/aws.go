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

func IsResourceAlreadyExists(err error) bool {
	return IsAWSAPIError(err, "ResourceAlreadyExistsException")
}

func IsEntitlementAlreadyExists(err error) bool {
	return IsAWSAPIError(err, "EntitlementAlreadyExistsException")
}

func IsAppStreamNotFound(err error) bool {
	return IsAWSAPIError(err, "ResourceNotFoundException", "EntitlementNotFoundException")
}
