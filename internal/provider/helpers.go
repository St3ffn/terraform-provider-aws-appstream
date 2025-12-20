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

func isAppStreamNotFound(err error) bool {
	var apiErr smithy.APIError
	if err == nil || !errors.As(err, &apiErr) {
		return false
	}

	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_ListEntitledApplications.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_AssociateApplicationToEntitlement.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DisassociateApplicationFromEntitlement.html
	// see https://docs.aws.amazon.com/appstream2/latest/APIReference/API_DescribeStacks.html
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

type assocDiagMode string

const (
	assocDiagPlan   assocDiagMode = "plan"
	assocDiagRead   assocDiagMode = "read"
	assocDiagDelete assocDiagMode = "delete"
)

func addAssocPartsDiagnostics(m associateApplicationEntitlementModel, diags *diag.Diagnostics, mode assocDiagMode) {
	if m.StackName.IsNull() || m.StackName.IsUnknown() ||
		m.EntitlementName.IsNull() || m.EntitlementName.IsUnknown() ||
		m.ApplicationIdentifier.IsNull() || m.ApplicationIdentifier.IsUnknown() {

		switch mode {
		case assocDiagPlan:
			diags.AddError(
				"Invalid Terraform Plan",
				"Cannot associate application to entitlement because stack_name, entitlement_name, and application_identifier must be known.",
			)
		case assocDiagDelete:
			diags.AddError(
				"Invalid Terraform State",
				"Cannot disassociate application from entitlement because stack_name, entitlement_name, and application_identifier must be known.",
			)
		case assocDiagRead:
			diags.AddError(
				"Invalid Terraform State",
				"Required attributes stack_name, entitlement_name, and application_identifier are missing from state. "+
					"This can happen after an incomplete import or a prior provider bug. Re-import or recreate the resource.",
			)
		}
	}
}

func boolOrNull(v *bool) types.Bool {
	if v == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*v)
}

func int32OrNull(v *int32) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*v))
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
