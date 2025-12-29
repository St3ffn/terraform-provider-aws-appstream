// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FlattenOwnedBool(prior types.Bool, awsValue *bool) types.Bool {
	// user never set it
	if prior.IsNull() {
		return types.BoolNull()
	}

	// terraform does not yet know
	if prior.IsUnknown() {
		return types.BoolUnknown()
	}

	// user owns it. reconcile
	if awsValue == nil {
		return types.BoolNull()
	}

	return types.BoolValue(*awsValue)
}

func FlattenOwnedInt32(prior types.Int32, awsValue *int32) types.Int32 {
	// user never set it
	if prior.IsNull() {
		return types.Int32Null()
	}

	// terraform does not yet know
	if prior.IsUnknown() {
		return types.Int32Unknown()
	}

	// user owns it. reconcile
	if awsValue == nil {
		return types.Int32Null()
	}

	return types.Int32Value(*awsValue)
}

func FlattenOwnedString(prior types.String, awsValue *string) types.String {
	// user never set it
	if prior.IsNull() {
		return types.StringNull()
	}

	// terraform does not yet know
	if prior.IsUnknown() {
		return types.StringUnknown()
	}

	// user owns it. reconcile
	if awsValue == nil {
		return types.StringNull()
	}

	return types.StringValue(*awsValue)
}

func FlattenOwnedStringSet(ctx context.Context, prior types.Set, awsValues []string, diags *diag.Diagnostics) types.Set {
	// user never managed this attribute
	if prior.IsNull() {
		return types.SetNull(types.StringType)
	}

	// terraform does not yet know during planning
	if prior.IsUnknown() {
		return types.SetUnknown(types.StringType)
	}

	// user owns this attribute; empty and nil mean empty
	if len(awsValues) == 0 {
		empty, d := types.SetValue(types.StringType, []attr.Value{})
		diags.Append(d...)
		return empty
	}

	setVal, d := types.SetValueFrom(ctx, types.StringType, awsValues)
	diags.Append(d...)
	if diags.HasError() {
		empty, _ := types.SetValue(types.StringType, []attr.Value{})
		return empty
	}

	return setVal
}
