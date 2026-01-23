// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestBoolPointerOrNil(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			got := BoolPointerOrNil(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %v", *got)
				}
				return
			}

			if got == nil || *got != *tt.want {
				t.Fatalf("BoolPointerOrNil(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestInt32PointerOrNil(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			got := Int32PointerOrNil(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %v", *got)
				}
				return
			}

			if got == nil || *got != *tt.want {
				t.Fatalf("Int32PointerOrNil(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestStringPointerOrNil(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			got := StringPointerOrNil(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %v", *got)
				}
				return
			}

			if got == nil || *got != *tt.want {
				t.Fatalf("StringPointerOrNil(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestExpandStringSetOrNil(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			var diags diag.Diagnostics

			got := ExpandStringSetOrNil(ctx, tt.input, &diags)

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ExpandStringSetOrNil() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestOptionalStringUpdate(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			var (
				called bool
				got    *string
			)

			setter := func(v *string) {
				called = true
				got = v
			}

			OptionalStringUpdate(tt.plan, tt.state, setter)

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

func TestDiffStringSets(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		oldSet  types.Set
		newSet  types.Set
		wantAdd []string
		wantRem []string
	}{
		{
			name:   "both_empty",
			oldSet: types.SetValueMust(types.StringType, []attr.Value{}),
			newSet: types.SetValueMust(types.StringType, []attr.Value{}),
		},
		{
			name: "only_additions",
			oldSet: types.SetValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("a")},
			),
			newSet: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("a"),
					types.StringValue("b"),
					types.StringValue("c"),
				},
			),
			wantAdd: []string{"b", "c"},
		},
		{
			name: "only_removals",
			oldSet: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("a"),
					types.StringValue("b"),
				},
			),
			newSet: types.SetValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("a")},
			),
			wantRem: []string{"b"},
		},
		{
			name: "add_and_remove",
			oldSet: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("a"),
					types.StringValue("b"),
				},
			),
			newSet: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("b"),
					types.StringValue("c"),
				},
			),
			wantAdd: []string{"c"},
			wantRem: []string{"a"},
		},
		{
			name: "identical_sets",
			oldSet: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("a"),
					types.StringValue("b"),
				},
			),
			newSet: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("b"),
					types.StringValue("a"),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			added, removed := DiffStringSets(ctx, tt.oldSet, tt.newSet, &diags)

			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			require.ElementsMatch(
				t,
				tt.wantAdd,
				added,
				"unexpected added elements",
			)

			require.ElementsMatch(
				t,
				tt.wantRem,
				removed,
				"unexpected removed elements",
			)
		})
	}
}
