// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func BoolOrNull(awsBool *bool) types.Bool {
	if awsBool == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*awsBool)
}

func Int32OrNull(awsInt32 *int32) types.Int32 {
	if awsInt32 == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*awsInt32)
}

func StringOrNull(awsString *string) types.String {
	if awsString == nil {
		return types.StringNull()
	}
	return types.StringValue(*awsString)
}

func StringFromTime(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(time.RFC3339))
}

func SetStringOrNull(ctx context.Context, values []string, diags *diag.Diagnostics) types.Set {
	// treat empty and nil slices the same here. this attribute is read-only and does not affect terraform behavior.
	if len(values) == 0 {
		return types.SetNull(types.StringType)
	}

	setVal, d := types.SetValueFrom(ctx, types.StringType, values)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(types.StringType)
	}

	return setVal
}

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
