// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BoolOrNull converts an AWS boolean pointer into a Terraform bool value.
func BoolOrNull(awsBool *bool) types.Bool {
	if awsBool == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*awsBool)
}

// Int32OrNull converts an AWS int32 pointer into a Terraform int32 value.
func Int32OrNull(awsInt32 *int32) types.Int32 {
	if awsInt32 == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*awsInt32)
}

// StringOrNull converts an AWS string pointer into a Terraform string value.
func StringOrNull(awsString *string) types.String {
	if awsString == nil {
		return types.StringNull()
	}
	return types.StringValue(*awsString)
}

// StringFromTime converts an AWS time pointer into an RFC3339 Terraform string.
func StringFromTime(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(time.RFC3339))
}

// SetStringOrNull converts a string slice into a Terraform set value.
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

// SetEnumStringOrNull converts a slice of AWS enum values into a Terraform
// set of strings and returns null for empty input.
func SetEnumStringOrNull[T ~string](ctx context.Context, values []T, diags *diag.Diagnostics) types.Set {
	// treat empty and nil slices the same here. this attribute is read-only and does not affect terraform behavior.
	if len(values) == 0 {
		return types.SetNull(types.StringType)
	}

	out := make([]string, 0, len(values))
	for _, v := range values {
		out = append(out, string(v))
	}

	setVal, d := types.SetValueFrom(ctx, types.StringType, out)
	diags.Append(d...)
	if diags.HasError() {
		return types.SetNull(types.StringType)
	}

	return setVal
}

// MapStringOrNull converts a Go string map into a Terraform map value.
func MapStringOrNull(ctx context.Context, m map[string]string, diags *diag.Diagnostics) types.Map {
	// treat empty and nil map the same here. this attribute is read-only and does not affect terraform behavior.
	if len(m) == 0 {
		return types.MapNull(types.StringType)
	}

	mapVal, d := types.MapValueFrom(ctx, types.StringType, m)
	diags.Append(d...)
	if diags.HasError() {
		return types.MapNull(types.StringType)
	}

	return mapVal
}
