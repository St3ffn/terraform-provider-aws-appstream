// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import "testing"

func TestIsRemoteField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   field
		want bool
	}{
		{
			name: "regular_optional_field",
			in: field{
				Tag:      "description",
				Optional: true,
			},
			want: true,
		},
		{
			name: "computed_only_field",
			in: field{
				Tag:      "created_time",
				Computed: true,
			},
			want: false,
		},
		{
			name: "tags_all_field",
			in: field{
				Tag: "tags_all",
			},
			want: false,
		},
		{
			name: "unannotated_field_default_true",
			in: field{
				Tag: "desired_state",
			},
			want: true,
		},
		{
			name: "remote_override_false",
			in: field{
				Tag:    "desired_state",
				Remote: boolPtr(false),
			},
			want: false,
		},
		{
			name: "remote_override_true",
			in: field{
				Tag:    "created_time",
				Remote: boolPtr(true),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isRemoteField(tt.in)
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func boolPtr(v bool) *bool {
	return &v
}

func TestOrChangedExpr_Empty(t *testing.T) {
	t.Parallel()

	got := orChangedExpr(nil, func(field) bool { return true }, "IsChanged")
	if got == nil {
		t.Fatalf("expected expression, got nil")
	}
}
