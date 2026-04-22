// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"context"
	"reflect"
	"testing"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExpandStates(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		input     types.Set
		want      map[awstypes.ImageState]struct{}
		wantError bool
	}{
		{
			name:  "null_set_returns_nil",
			input: types.SetNull(types.StringType),
			want:  nil,
		},
		{
			name:  "unknown_set_returns_nil",
			input: types.SetUnknown(types.StringType),
			want:  nil,
		},
		{
			name:  "empty_set_returns_nil",
			input: types.SetValueMust(types.StringType, []attr.Value{}),
			want:  nil,
		},
		{
			name: "single_value_returns_map",
			input: types.SetValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("AVAILABLE")},
			),
			want: map[awstypes.ImageState]struct{}{
				awstypes.ImageStateAvailable: {},
			},
		},
		{
			name: "multiple_values_returns_map",
			input: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("PENDING"),
					types.StringValue("VALIDATING"),
				},
			),
			want: map[awstypes.ImageState]struct{}{
				awstypes.ImageStatePending:    {},
				awstypes.ImageStateValidating: {},
			},
		},
		{
			name: "invalid_value_produces_diagnostics",
			input: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringUnknown(),
				},
			),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			got := expandStates(ctx, tt.input, &diags)

			if tt.wantError {
				if !diags.HasError() {
					t.Fatalf("expected diagnostics error, got none")
				}
				return
			}

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expandStates() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
