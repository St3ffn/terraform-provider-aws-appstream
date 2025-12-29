// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func TestARNValidator_Constructors(t *testing.T) {
	tests := []struct {
		name      string
		validator validator.String
		value     basetypes.StringValue
		wantError bool
	}{
		{
			name:      "null_value_returns_no_error",
			validator: ValidARN(),
			value:     types.StringNull(),
			wantError: false,
		},
		{
			name:      "unknown_value_returns_no_error",
			validator: ValidARN(),
			value:     types.StringUnknown(),
			wantError: false,
		},
		{
			name:      "valid_arn_returns_no_error",
			validator: ValidARN(),
			value:     types.StringValue("arn:aws:iam::123456789012:role/MyRole"),
			wantError: false,
		},
		{
			name:      "invalid_arn_returns_error",
			validator: ValidARN(),
			value:     types.StringValue("not-an-arn"),
			wantError: true,
		},
		{
			name:      "valid_arn_with_service_returns_no_error",
			validator: ValidARNWithService("iam"),
			value:     types.StringValue("arn:aws:iam::123456789012:role/MyRole"),
			wantError: false,
		},
		{
			name:      "wrong_service_returns_error",
			validator: ValidARNWithService("iam"),
			value:     types.StringValue("arn:aws:appstream:us-east-1:123456789012:other/my-app"),
			wantError: true,
		},
		{
			name:      "valid_arn_with_service_and_resource_returns_no_error",
			validator: ValidARNWithServiceAndResource("appstream", "application/"),
			value:     types.StringValue("arn:aws:appstream:us-east-1:123456789012:application/my-app"),
			wantError: false,
		},
		{
			name:      "wrong_resource_returns_error",
			validator: ValidARNWithServiceAndResource("appstream", "application/"),
			value:     types.StringValue("arn:aws:appstream:us-east-1:123456789012:other/my-image"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: tt.value,
			}
			resp := &validator.StringResponse{}

			tt.validator.ValidateString(context.Background(), req, resp)

			if tt.wantError {
				require.True(
					t,
					resp.Diagnostics.HasError(),
					"expected validation error but got none",
				)
			} else {
				require.False(
					t,
					resp.Diagnostics.HasError(),
					"unexpected validation error: %v",
					resp.Diagnostics,
				)
			}
		})
	}
}
