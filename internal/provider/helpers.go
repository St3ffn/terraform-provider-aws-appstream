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

func isContextCanceled(ctx context.Context) bool {
	return errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded)
}

func isAppStreamNotFound(err error) bool {
	var apiErr smithy.APIError
	if err == nil || !errors.As(err, &apiErr) {
		return false
	}

	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_ListEntitledApplications.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_AssociateApplicationToEntitlement.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DisassociateApplicationFromEntitlement.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeStacks.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeFleets.html
	switch apiErr.ErrorCode() {
	case "ResourceNotFoundException", "EntitlementNotFoundException":
		return true
	default:
		return false
	}
}

func isAppStreamAlreadyExists(err error) bool {
	var apiErr smithy.APIError
	if err == nil || !errors.As(err, &apiErr) {
		return false
	}
	return apiErr.ErrorCode() == "ResourceAlreadyExistsException"
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
