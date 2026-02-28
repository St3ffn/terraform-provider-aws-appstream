// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package fleet

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDesiredStateFromPlan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   types.String
		want desiredState
	}{
		{
			name: "null_defaults_to_inherit",
			in:   types.StringNull(),
			want: desiredStateInherit,
		},
		{
			name: "unknown_defaults_to_inherit",
			in:   types.StringUnknown(),
			want: desiredStateInherit,
		},
		{
			name: "running",
			in:   types.StringValue(desiredStateRunning.String()),
			want: desiredStateRunning,
		},
		{
			name: "inherit",
			in:   types.StringValue(desiredStateInherit.String()),
			want: desiredStateInherit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := desiredStateFromPlan(tt.in)
			if got != tt.want {
				t.Fatalf("got desired state %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUpdateBehaviorFromPlan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   types.String
		want updateBehavior
	}{
		{
			name: "null_defaults_to_auto",
			in:   types.StringNull(),
			want: updateBehaviorAutoStopStart,
		},
		{
			name: "unknown_defaults_to_auto",
			in:   types.StringUnknown(),
			want: updateBehaviorAutoStopStart,
		},
		{
			name: "configured",
			in:   types.StringValue(updateBehaviorFailIfRunning.String()),
			want: updateBehaviorFailIfRunning,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := updateBehaviorFromPlan(tt.in)
			if got != tt.want {
				t.Fatalf("got update behavior %q, want %q", got, tt.want)
			}
		})
	}
}
