// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import "testing"

func TestParseRemoteOverride(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want bool
		ok   bool
	}{
		{
			name: "true",
			text: "codegen:has_remote_changes=true",
			want: true,
			ok:   true,
		},
		{
			name: "false",
			text: "codegen:has_remote_changes=false",
			want: false,
			ok:   true,
		},
		{
			name: "with_spaces",
			text: "codegen:has_remote_changes = false",
			want: false,
			ok:   true,
		},
		{
			name: "no_annotation",
			text: "something else",
			want: false,
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := parseRemoteOverride(tt.text)
			if ok != tt.ok {
				t.Fatalf("ok=%v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
