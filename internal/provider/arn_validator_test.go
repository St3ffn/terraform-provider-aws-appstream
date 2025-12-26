// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestARNValidator_Constructors(t *testing.T) {
	tests := []struct {
		name      string
		validator validator.String
		value     basetypes.StringValue
		wantError bool
	}{
		{
			name:      "validARN - null value - success",
			validator: validARN(),
			value:     types.StringNull(),
			wantError: false,
		},
		{
			name:      "validARN - unknown value - success",
			validator: validARN(),
			value:     types.StringUnknown(),
			wantError: false,
		},
		{
			name:      "validARN - success",
			validator: validARN(),
			value:     types.StringValue("arn:aws:iam::123456789012:role/MyRole"),
			wantError: false,
		},
		{
			name:      "validARN - invalid",
			validator: validARN(),
			value:     types.StringValue("not-an-arn"),
			wantError: true,
		},
		{
			name:      "validARNWithService - success",
			validator: validARNWithService("iam"),
			value:     types.StringValue("arn:aws:iam::123456789012:role/MyRole"),
			wantError: false,
		},
		{
			name:      "validARNWithService - wrong service",
			validator: validARNWithService("iam"),
			value:     types.StringValue("arn:aws:appstream:us-east-1:123456789012:other/my-app"),
			wantError: true,
		},
		{
			name:      "validARNWithServiceAndResource - success",
			validator: validARNWithServiceAndResource("appstream", "application/"),
			value:     types.StringValue("arn:aws:appstream:us-east-1:123456789012:application/my-app"),
			wantError: false,
		},
		{
			name:      "validARNWithServiceAndResource - wrong resource",
			validator: validARNWithServiceAndResource("appstream", "application/"),
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

			if tt.wantError && !resp.Diagnostics.HasError() {
				t.Fatalf("expected validation error, got none")
			}

			if !tt.wantError && resp.Diagnostics.HasError() {
				t.Fatalf("unexpected validation error: %v", resp.Diagnostics)
			}
		})
	}
}
