// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"errors"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/stretchr/testify/require"
)

func TestErrorPredicates(t *testing.T) {
	type testCase struct {
		name string
		err  error
		fn   func(error) bool
		want bool
	}

	tests := []testCase{
		{
			name: "operation_not_permitted/match",
			err:  &smithy.GenericAPIError{Code: "OperationNotPermittedException"},
			fn:   IsOperationNotPermittedException,
			want: true,
		},
		{
			name: "operation_not_permitted/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   IsOperationNotPermittedException,
			want: false,
		},
		{
			name: "operation_not_permitted/no_match_nil",
			err:  nil,
			fn:   IsOperationNotPermittedException,
			want: false,
		},
		{
			name: "resource_not_found/match",
			err:  &smithy.GenericAPIError{Code: "ResourceNotFoundException"},
			fn:   IsResourceNotFoundException,
			want: true,
		},
		{
			name: "resource_not_found/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   IsResourceNotFoundException,
			want: false,
		},
		{
			name: "resource_not_found/no_match_nil",
			err:  nil,
			fn:   IsResourceNotFoundException,
			want: false,
		},
		{
			name: "concurrent_modification/match",
			err:  &smithy.GenericAPIError{Code: "ConcurrentModificationException"},
			fn:   IsConcurrentModificationException,
			want: true,
		},
		{
			name: "concurrent_modification/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   IsConcurrentModificationException,
			want: false,
		},
		{
			name: "concurrent_modification/no_match_nil",
			err:  nil,
			fn:   IsConcurrentModificationException,
			want: false,
		},
		{
			name: "entitlement_not_found/match",
			err:  &smithy.GenericAPIError{Code: "EntitlementNotFoundException"},
			fn:   IsEntitlementNotFoundException,
			want: true,
		},
		{
			name: "entitlement_not_found/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   IsEntitlementNotFoundException,
			want: false,
		},
		{
			name: "entitlement_not_found/no_match_nil",
			err:  nil,
			fn:   IsEntitlementNotFoundException,
			want: false,
		},
		{
			name: "resource_not_available/match",
			err:  &smithy.GenericAPIError{Code: "ResourceNotAvailableException"},
			fn:   IsResourceNotAvailableException,
			want: true,
		},
		{
			name: "resource_not_available/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   IsResourceNotAvailableException,
			want: false,
		},
		{
			name: "resource_not_available/no_match_nil",
			err:  nil,
			fn:   IsResourceNotAvailableException,
			want: false,
		},
		{
			name: "resource_already_exists/match",
			err:  &smithy.GenericAPIError{Code: "ResourceAlreadyExistsException"},
			fn:   IsResourceAlreadyExists,
			want: true,
		},
		{
			name: "resource_already_exists/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   IsResourceAlreadyExists,
			want: false,
		},
		{
			name: "resource_already_exists/no_match_nil",
			err:  nil,
			fn:   IsResourceAlreadyExists,
			want: false,
		},
		{
			name: "entitlement_already_exists/match",
			err:  &smithy.GenericAPIError{Code: "EntitlementAlreadyExistsException"},
			fn:   IsEntitlementAlreadyExists,
			want: true,
		},
		{
			name: "entitlement_already_exists/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   IsEntitlementAlreadyExists,
			want: false,
		},
		{
			name: "entitlement_already_exists/no_match_nil",
			err:  nil,
			fn:   IsEntitlementAlreadyExists,
			want: false,
		},
		{
			name: "appstream_not_found/resource_not_found",
			err:  &smithy.GenericAPIError{Code: "ResourceNotFoundException"},
			fn:   IsAppStreamNotFound,
			want: true,
		},
		{
			name: "appstream_not_found/entitlement_not_found",
			err:  &smithy.GenericAPIError{Code: "EntitlementNotFoundException"},
			fn:   IsAppStreamNotFound,
			want: true,
		},
		{
			name: "appstream_not_found/other_error",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   IsAppStreamNotFound,
			want: false,
		},
		{
			name: "appstream_not_found/no_match_nil",
			err:  nil,
			fn:   IsAppStreamNotFound,
			want: false,
		},
		{
			name: "non_aws_error",
			err:  errors.New("plain error"),
			fn:   IsResourceNotFoundException,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.err)
			require.Equal(t, tt.want, got)
		})
	}
}
