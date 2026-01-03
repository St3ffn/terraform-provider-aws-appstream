// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import "github.com/hashicorp/terraform-plugin-framework/types"

func FlattenOptionalComputedBool(prior types.Bool, awsValue *bool) types.Bool {
	if awsValue != nil {
		return types.BoolValue(*awsValue)
	}

	// if aws didn't return it, keep prior
	if prior.IsUnknown() {
		return types.BoolUnknown()
	}
	if prior.IsNull() {
		return types.BoolNull()
	}
	return types.BoolValue(prior.ValueBool())
}

func FlattenOptionalComputedString(prior types.String, awsValue *string) types.String {
	if awsValue != nil {
		return types.StringValue(*awsValue)
	}

	// if aws didn't return it, keep prior
	if prior.IsUnknown() {
		return types.StringUnknown()
	}
	if prior.IsNull() {
		return types.StringNull()
	}
	return types.StringValue(prior.ValueString())
}
