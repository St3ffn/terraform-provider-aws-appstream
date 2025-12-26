// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func isContextCanceled(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func isAWSAPIError(err error, code ...string) bool {
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

func isOperationNotPermittedException(err error) bool {
	return isAWSAPIError(err, "OperationNotPermittedException")
}

func isResourceNotFoundException(err error) bool {
	return isAWSAPIError(err, "ResourceNotFoundException")
}

func isConcurrentModificationException(err error) bool {
	return isAWSAPIError(err, "ConcurrentModificationException")
}

func isEntitlementNotFoundException(err error) bool {
	return isAWSAPIError(err, "EntitlementNotFoundException")
}

func isResourceAlreadyExists(err error) bool {
	return isAWSAPIError(err, "ResourceAlreadyExistsException")
}

func isAppStreamNotFound(err error) bool {
	return isAWSAPIError(err, "ResourceNotFoundException", "EntitlementNotFoundException")
}

func boolOrNull(v *bool) types.Bool {
	if v == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*v)
}

func int32OrNull(v *int32) types.Int32 {
	if v == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*v)
}

func stringOrNull(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

func stringFromTime(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(time.RFC3339))
}

func setStringOrNull(
	ctx context.Context,
	values []string,
	diags *diag.Diagnostics,
) types.Set {

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

func boolPointerOrNil(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

func int32PointerOrNil(v types.Int32) *int32 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	i := v.ValueInt32()
	return &i
}

func stringPointerOrNil(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}

func expandStringSetOrNil(ctx context.Context, set types.Set, diags *diag.Diagnostics) []string {
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

func optionalStringUpdate(
	plan types.String,
	state types.String,
	setter func(*string),
) {
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
