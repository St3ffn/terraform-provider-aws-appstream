// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestChanged_string(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		state  types.String
		plan   types.String
		expect bool
	}{
		{"plan unknown", types.StringValue("a"), types.StringUnknown(), false},
		{"both null", types.StringNull(), types.StringNull(), false},
		{"null to value", types.StringNull(), types.StringValue("a"), true},
		{"value to null", types.StringValue("a"), types.StringNull(), true},
		{"same value", types.StringValue("a"), types.StringValue("a"), false},
		{"different value", types.StringValue("a"), types.StringValue("b"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expect, Changed(tt.state, tt.plan))
		})
	}
}

func TestChanged_bool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		state  types.Bool
		plan   types.Bool
		expect bool
	}{
		{"plan unknown", types.BoolValue(true), types.BoolUnknown(), false},
		{"both null", types.BoolNull(), types.BoolNull(), false},
		{"null to value", types.BoolNull(), types.BoolValue(true), true},
		{"value to null", types.BoolValue(true), types.BoolNull(), true},
		{"same value", types.BoolValue(true), types.BoolValue(true), false},
		{"different value", types.BoolValue(true), types.BoolValue(false), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expect, Changed(tt.state, tt.plan))
		})
	}
}

func TestChanged_int32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		state  types.Int32
		plan   types.Int32
		expect bool
	}{
		{"plan unknown", types.Int32Value(1), types.Int32Unknown(), false},
		{"both null", types.Int32Null(), types.Int32Null(), false},
		{"null to value", types.Int32Null(), types.Int32Value(1), true},
		{"value to null", types.Int32Value(1), types.Int32Null(), true},
		{"same value", types.Int32Value(1), types.Int32Value(1), false},
		{"different value", types.Int32Value(1), types.Int32Value(2), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expect, Changed(tt.state, tt.plan))
		})
	}
}

func TestChanged_set(t *testing.T) {
	t.Parallel()

	objectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		},
	}

	mustObject := func(t *testing.T, id, name string) attr.Value {
		t.Helper()

		obj, diags := types.ObjectValue(
			objectType.AttrTypes,
			map[string]attr.Value{
				"id":   types.StringValue(id),
				"name": types.StringValue(name),
			},
		)
		require.False(t, diags.HasError())
		return obj
	}

	tests := []struct {
		name   string
		state  types.Set
		plan   types.Set
		expect bool
	}{
		{
			name:   "plan unknown",
			state:  types.SetNull(objectType),
			plan:   types.SetUnknown(objectType),
			expect: false,
		},
		{
			name:   "both null",
			state:  types.SetNull(objectType),
			plan:   types.SetNull(objectType),
			expect: false,
		},
		{
			name:   "null to empty",
			state:  types.SetNull(objectType),
			plan:   types.SetValueMust(objectType, []attr.Value{}),
			expect: true,
		},
		{
			name: "same objects different order",
			state: types.SetValueMust(objectType, []attr.Value{
				mustObject(t, "1", "a"),
				mustObject(t, "2", "b"),
			}),
			plan: types.SetValueMust(objectType, []attr.Value{
				mustObject(t, "2", "b"),
				mustObject(t, "1", "a"),
			}),
			expect: false,
		},
		{
			name: "nested field changed",
			state: types.SetValueMust(objectType, []attr.Value{
				mustObject(t, "1", "a"),
			}),
			plan: types.SetValueMust(objectType, []attr.Value{
				mustObject(t, "1", "changed"),
			}),
			expect: true,
		},
		{
			name: "object added",
			state: types.SetValueMust(objectType, []attr.Value{
				mustObject(t, "1", "a"),
			}),
			plan: types.SetValueMust(objectType, []attr.Value{
				mustObject(t, "1", "a"),
				mustObject(t, "2", "b"),
			}),
			expect: true,
		},
		{
			name: "object removed",
			state: types.SetValueMust(objectType, []attr.Value{
				mustObject(t, "1", "a"),
				mustObject(t, "2", "b"),
			}),
			plan: types.SetValueMust(objectType, []attr.Value{
				mustObject(t, "1", "a"),
			}),
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			changed := Changed(tt.state, tt.plan)
			require.Equal(t, tt.expect, changed)
		})
	}
}

func TestChanged_object(t *testing.T) {
	t.Parallel()

	attrTypes := map[string]attr.Type{
		"foo": types.StringType,
	}

	mustObjectValue := func(t *testing.T, value string) types.Object {
		t.Helper()

		obj, diags := types.ObjectValue(
			attrTypes,
			map[string]attr.Value{
				"foo": types.StringValue(value),
			},
		)
		require.False(t, diags.HasError())

		return obj
	}

	tests := []struct {
		name   string
		state  types.Object
		plan   types.Object
		expect bool
	}{
		{
			name:   "plan unknown",
			state:  types.ObjectNull(attrTypes),
			plan:   types.ObjectUnknown(attrTypes),
			expect: false,
		},
		{
			name:   "state null plan non-null",
			state:  types.ObjectNull(attrTypes),
			plan:   mustObjectValue(t, "bar"),
			expect: true,
		},
		{
			name:   "state non-null plan null",
			state:  mustObjectValue(t, "bar"),
			plan:   types.ObjectNull(attrTypes),
			expect: true,
		},
		{
			name:   "both null",
			state:  types.ObjectNull(attrTypes),
			plan:   types.ObjectNull(attrTypes),
			expect: false,
		},
		{
			name:   "both non-null same value",
			state:  mustObjectValue(t, "bar"),
			plan:   mustObjectValue(t, "bar"),
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			changed := Changed(tt.state, tt.plan)
			require.Equal(t, tt.expect, changed)
		})
	}
}
