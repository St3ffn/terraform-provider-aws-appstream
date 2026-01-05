// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
			got := BoolOrNull(tt.input)

			if !got.Equal(tt.want) {
				t.Fatalf("BoolOrNull(%v) = %v, want %v", tt.input, got, tt.want)
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
			got := Int32OrNull(tt.input)

			if !got.Equal(tt.want) {
				t.Fatalf("Int32OrNull(%v) = %v, want %v", tt.input, got, tt.want)
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
			got := StringOrNull(tt.input)

			if !got.Equal(tt.want) {
				t.Fatalf("StringOrNull(%v) = %v, want %v", tt.input, got, tt.want)
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
			got := StringFromTime(tt.input)

			if !got.Equal(tt.want) {
				t.Fatalf("StringFromTime(%v) = %v, want %v", tt.input, got, tt.want)
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

			got := SetStringOrNull(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("TestSetStringOrNull(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSetEnumStringOrNull(t *testing.T) {
	ctx := context.Background()

	type testEnum string

	const (
		TestEnumOne testEnum = "ONE"
		TestEnumTwo testEnum = "TWO"
	)

	tests := []struct {
		name  string
		input []testEnum
		want  types.Set
	}{
		{
			name:  "nil_slice_returns_null",
			input: nil,
			want:  types.SetNull(types.StringType),
		},
		{
			name:  "empty_slice_returns_null",
			input: []testEnum{},
			want:  types.SetNull(types.StringType),
		},
		{
			name:  "single_value_returns_set",
			input: []testEnum{TestEnumOne},
			want: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("ONE"),
				},
			),
		},
		{
			name:  "multiple_values_returns_set",
			input: []testEnum{TestEnumOne, TestEnumTwo},
			want: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("ONE"),
					types.StringValue("TWO"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := SetEnumStringOrNull(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf(
					"SetEnumStringOrNull(%v) = %v, want %v",
					tt.input,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestMapStringOrNull(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input map[string]string
		want  types.Map
	}{
		{
			name:  "nil_map_returns_null",
			input: nil,
			want:  types.MapNull(types.StringType),
		},
		{
			name:  "empty_map_returns_null",
			input: map[string]string{},
			want:  types.MapNull(types.StringType),
		},
		{
			name: "single_entry_returns_map",
			input: map[string]string{
				"key": "value",
			},
			want: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key": types.StringValue("value"),
				},
			),
		},
		{
			name: "multiple_entries_returns_map",
			input: map[string]string{
				"one": "1",
				"two": "2",
			},
			want: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one": types.StringValue("1"),
					"two": types.StringValue("2"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics

			got := MapStringOrNull(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("MapStringOrNull(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
