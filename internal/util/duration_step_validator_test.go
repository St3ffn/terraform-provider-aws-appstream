// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestDurationStepValidator_ValidateInt32(t *testing.T) {
	ctx := context.Background()

	v := DurationWithStep(
		60,    // min
		36000, // max
		60,    // step
		true,  // allowZero
	)

	type testCase struct {
		name          string
		value         types.Int32
		expectError   bool
		errorContains string
	}

	tests := []testCase{
		{
			name:        "null value is allowed",
			value:       types.Int32Null(),
			expectError: false,
		},
		{
			name:        "unknown value is allowed",
			value:       types.Int32Unknown(),
			expectError: false,
		},
		{
			name:        "zero is allowed when allowZero is true",
			value:       types.Int32Value(0),
			expectError: false,
		},
		{
			name:        "minimum valid value",
			value:       types.Int32Value(60),
			expectError: false,
		},
		{
			name:        "maximum valid value",
			value:       types.Int32Value(36000),
			expectError: false,
		},
		{
			name:        "valid multiple of step",
			value:       types.Int32Value(120),
			expectError: false,
		},
		{
			name:          "value below minimum is rejected",
			value:         types.Int32Value(30),
			expectError:   true,
			errorContains: "between 60 and 36000",
		},
		{
			name:          "value above maximum is rejected",
			value:         types.Int32Value(40000),
			expectError:   true,
			errorContains: "between 60 and 36000",
		},
		{
			name:          "value not divisible by step is rejected",
			value:         types.Int32Value(90),
			expectError:   true,
			errorContains: "increments of 60",
		},
		{
			name:          "non-zero invalid step below min still reports min/max",
			value:         types.Int32Value(1),
			expectError:   true,
			errorContains: "between 60 and 36000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := validator.Int32Request{
				ConfigValue: tc.value,
				Path:        path.Root("idle_disconnect_timeout_in_seconds"),
			}

			resp := &validator.Int32Response{}

			v.ValidateInt32(ctx, req, resp)

			if tc.expectError {
				require.True(
					t,
					resp.Diagnostics.HasError(),
					"expected validation error for test case %q, but got none",
					tc.name,
				)

				if tc.errorContains != "" {
					found := false
					for _, d := range resp.Diagnostics {
						if d.Detail() != "" && contains(d.Detail(), tc.errorContains) {
							found = true
							break
						}
					}

					require.True(
						t,
						found,
						"expected diagnostic for test case %q to contain %q, but diagnostics were: %+v",
						tc.name,
						tc.errorContains,
						resp.Diagnostics,
					)
				}
			} else {
				require.False(
					t,
					resp.Diagnostics.HasError(),
					"did not expect validation error for test case %q, but got: %+v",
					tc.name,
					resp.Diagnostics,
				)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && fmt.Sprintf("%v", s) != "" && (func() bool {
		return stringContains(s, substr)
	})())
}

func stringContains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
