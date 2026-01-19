// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func BoolPointerOrNil(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

func Int32PointerOrNil(v types.Int32) *int32 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	i := v.ValueInt32()
	return &i
}

func StringPointerOrNil(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}

func ExpandStringSetOrNil(ctx context.Context, set types.Set, diags *diag.Diagnostics) []string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	var values []string
	diags.Append(set.ElementsAs(ctx, &values, false)...)
	if diags.HasError() {
		return nil
	}

	if len(values) == 0 {
		return nil
	}

	return values
}

func OptionalStringUpdate(plan types.String, state types.String, setter func(*string)) {
	switch {
	case plan.IsUnknown():
		return
	case !plan.IsNull():
		v := plan.ValueString()
		setter(&v)
	case plan.IsNull() && !state.IsNull():
		empty := ""
		setter(&empty)
	}
}

func DiffStringSets(ctx context.Context, oldSet, newSet types.Set, diags *diag.Diagnostics) (added, removed []string) {
	var oldVals, newVals []string

	diags.Append(oldSet.ElementsAs(ctx, &oldVals, false)...)
	diags.Append(newSet.ElementsAs(ctx, &newVals, false)...)
	if diags.HasError() {
		return
	}

	oldIndex := make(map[string]struct{}, len(oldVals))
	for _, v := range oldVals {
		oldIndex[v] = struct{}{}
	}

	newIndex := make(map[string]struct{}, len(newVals))
	for _, v := range newVals {
		newIndex[v] = struct{}{}
	}

	for v := range newIndex {
		if _, ok := oldIndex[v]; !ok {
			added = append(added, v)
		}
	}

	for v := range oldIndex {
		if _, ok := newIndex[v]; !ok {
			removed = append(removed, v)
		}
	}

	return
}
