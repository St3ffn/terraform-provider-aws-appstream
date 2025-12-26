// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			name: "context_canceled/matches_canceled",
			err:  context.Canceled,
			fn:   isContextCanceled,
			want: true,
		},
		{
			name: "context_canceled/matches_deadline_exceeded",
			err:  context.DeadlineExceeded,
			fn:   isContextCanceled,
			want: true,
		},
		{
			name: "context_canceled/other_error_no_match",
			err:  errors.New("other"),
			fn:   isContextCanceled,
			want: false,
		},
		{
			name: "operation_not_permitted/match",
			err:  &smithy.GenericAPIError{Code: "OperationNotPermittedException"},
			fn:   isOperationNotPermittedException,
			want: true,
		},
		{
			name: "operation_not_permitted/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   isOperationNotPermittedException,
			want: false,
		},
		{
			name: "resource_not_found/match",
			err:  &smithy.GenericAPIError{Code: "ResourceNotFoundException"},
			fn:   isResourceNotFoundException,
			want: true,
		},
		{
			name: "resource_not_found/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   isResourceNotFoundException,
			want: false,
		},
		{
			name: "concurrent_modification/match",
			err:  &smithy.GenericAPIError{Code: "ConcurrentModificationException"},
			fn:   isConcurrentModificationException,
			want: true,
		},
		{
			name: "concurrent_modification/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   isConcurrentModificationException,
			want: false,
		},
		{
			name: "entitlement_not_found/match",
			err:  &smithy.GenericAPIError{Code: "EntitlementNotFoundException"},
			fn:   isEntitlementNotFoundException,
			want: true,
		},
		{
			name: "entitlement_not_found/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   isEntitlementNotFoundException,
			want: false,
		},
		{
			name: "resource_already_exists/match",
			err:  &smithy.GenericAPIError{Code: "ResourceAlreadyExistsException"},
			fn:   isResourceAlreadyExists,
			want: true,
		},
		{
			name: "resource_already_exists/no_match",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   isResourceAlreadyExists,
			want: false,
		},
		{
			name: "appstream_not_found/resource_not_found",
			err:  &smithy.GenericAPIError{Code: "ResourceNotFoundException"},
			fn:   isAppStreamNotFound,
			want: true,
		},
		{
			name: "appstream_not_found/entitlement_not_found",
			err:  &smithy.GenericAPIError{Code: "EntitlementNotFoundException"},
			fn:   isAppStreamNotFound,
			want: true,
		},
		{
			name: "appstream_not_found/other_error",
			err:  &smithy.GenericAPIError{Code: "other"},
			fn:   isAppStreamNotFound,
			want: false,
		},
		{
			name: "non_aws_error",
			err:  errors.New("plain error"),
			fn:   isResourceNotFoundException,
			want: false,
		},
		{
			name: "nil_error",
			err:  nil,
			fn:   isResourceNotFoundException,
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

func TestBoolOrNull(t *testing.T) {
	tests := []struct {
		name  string
		input *bool
		want  types.Bool
	}{
		{name: "nil_return_null", input: nil, want: types.BoolNull()},
		{name: "true_returns_true", input: aws.Bool(true), want: types.BoolValue(true)},
		{name: "false_returns_false", input: aws.Bool(false), want: types.BoolValue(false)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := boolOrNull(tt.input)

			if !got.Equal(tt.want) {
				t.Fatalf("boolOrNull(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestInt32OrNull(t *testing.T) {
	tests := []struct {
		name  string
		input *int32
		want  types.Int32
	}{
		{name: "nil_returns_null", input: nil, want: types.Int32Null()},
		{name: "zero_returns_zero", input: aws.Int32(0), want: types.Int32Value(0)},
		{name: "positive_returns_value", input: aws.Int32(42), want: types.Int32Value(42)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := int32OrNull(tt.input)

			if !got.Equal(tt.want) {
				t.Fatalf("int32OrNull(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestStringOrNull(t *testing.T) {
	tests := []struct {
		name  string
		input *string
		want  types.String
	}{
		{name: "nil_returns_null", input: nil, want: types.StringNull()},
		{name: "empty_string_returns_empty", input: aws.String(""), want: types.StringValue("")},
		{name: "value_returns_value", input: aws.String("hello"), want: types.StringValue("hello")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringOrNull(tt.input)

			if !got.Equal(tt.want) {
				t.Fatalf("stringOrNull(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestStringFromTime(t *testing.T) {
	t1 := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name  string
		input *time.Time
		want  types.String
	}{
		{name: "nil_returns_null", input: nil, want: types.StringNull()},
		{name: "time_returns_rfc3339", input: &t1, want: types.StringValue("2024-01-02T03:04:05Z")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringFromTime(tt.input)

			if !got.Equal(tt.want) {
				t.Fatalf("stringFromTime(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSetStringOrNull(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input []string
		want  types.Set
	}{
		{
			name:  "nil_slice_returns_null",
			input: nil,
			want:  types.SetNull(types.StringType),
		},
		{
			name:  "empty_slice_returns_null",
			input: []string{},
			want:  types.SetNull(types.StringType),
		},
		{
			name:  "single_value_returns_set",
			input: []string{"one"},
			want: types.SetValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("one")},
			),
		},
		{
			name:  "multiple_values_returns_set",
			input: []string{"one", "two"},
			want: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("one"),
					types.StringValue("two"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := setStringOrNull(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("setStringOrNull(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestBoolPointerOrNil(t *testing.T) {
	tests := []struct {
		name  string
		input types.Bool
		want  *bool
	}{
		{name: "null_returns_nil", input: types.BoolNull(), want: nil},
		{name: "unknown_returns_nil", input: types.BoolUnknown(), want: nil},
		{name: "true_returns_pointer", input: types.BoolValue(true), want: aws.Bool(true)},
		{name: "false_returns_pointer", input: types.BoolValue(false), want: aws.Bool(false)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := boolPointerOrNil(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %v", *got)
				}
				return
			}

			if got == nil || *got != *tt.want {
				t.Fatalf("boolPointerOrNil(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestInt32PointerOrNil(t *testing.T) {
	tests := []struct {
		name  string
		input types.Int32
		want  *int32
	}{
		{name: "null_returns_nil", input: types.Int32Null(), want: nil},
		{name: "unknown_returns_nil", input: types.Int32Unknown(), want: nil},
		{name: "value_returns_pointer", input: types.Int32Value(42), want: aws.Int32(42)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := int32PointerOrNil(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %v", *got)
				}
				return
			}

			if got == nil || *got != *tt.want {
				t.Fatalf("int32PointerOrNil(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestStringPointerOrNil(t *testing.T) {
	tests := []struct {
		name  string
		input types.String
		want  *string
	}{
		{name: "null_returns_nil", input: types.StringNull(), want: nil},
		{name: "unknown_returns_nil", input: types.StringUnknown(), want: nil},
		{name: "value_returns_pointer", input: types.StringValue("hello"), want: aws.String("hello")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringPointerOrNil(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %v", *got)
				}
				return
			}

			if got == nil || *got != *tt.want {
				t.Fatalf("stringPointerOrNil(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestExpandStringSetOrNil(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     types.Set
		want      []string
		wantError bool
	}{
		{
			name:  "null_set_returns_nil",
			input: types.SetNull(types.StringType),
			want:  nil,
		},
		{
			name:  "unknown_set_returns_nil",
			input: types.SetUnknown(types.StringType),
			want:  nil,
		},
		{
			name:  "empty_set_returns_nil",
			input: types.SetValueMust(types.StringType, []attr.Value{}),
			want:  nil,
		},
		{
			name: "single_value_returns_slice",
			input: types.SetValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("one")},
			),
			want: []string{"one"},
		},
		{
			name: "multiple_values_returns_slice",
			input: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("one"),
					types.StringValue("two"),
				},
			),
			want: []string{"one", "two"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := expandStringSetOrNil(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expandStringSetOrNil() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestOptionalStringUpdate(t *testing.T) {
	tests := []struct {
		name       string
		plan       types.String
		state      types.String
		wantCalled bool
		wantValue  *string
	}{
		{
			name:       "plan_unknown_noop",
			plan:       types.StringUnknown(),
			state:      types.StringValue("old"),
			wantCalled: false,
		},
		{
			name:       "plan_value_sets_value",
			plan:       types.StringValue("new"),
			state:      types.StringValue("old"),
			wantCalled: true,
			wantValue:  aws.String("new"),
		},
		{
			name:       "plan_value_state_null_sets_value",
			plan:       types.StringValue("new"),
			state:      types.StringNull(),
			wantCalled: true,
			wantValue:  aws.String("new"),
		},
		{
			name:       "plan_null_state_value_clears",
			plan:       types.StringNull(),
			state:      types.StringValue("old"),
			wantCalled: true,
			wantValue:  aws.String(""),
		},
		{
			name:       "plan_null_state_null_noop",
			plan:       types.StringNull(),
			state:      types.StringNull(),
			wantCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				called bool
				got    *string
			)

			setter := func(v *string) {
				called = true
				got = v
			}

			optionalStringUpdate(tt.plan, tt.state, setter)

			if called != tt.wantCalled {
				t.Fatalf("setter called = %v, want %v", called, tt.wantCalled)
			}

			if !tt.wantCalled {
				return
			}

			if tt.wantValue == nil {
				if got != nil {
					t.Fatalf("expected nil value, got %v", *got)
				}
				return
			}

			if got == nil || *got != *tt.wantValue {
				t.Fatalf("got %v, want %v", got, tt.wantValue)
			}
		})
	}
}
