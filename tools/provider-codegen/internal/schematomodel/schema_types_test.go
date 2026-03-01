// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import "testing"

func TestSchemaTypeToModelType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		kind string
		want string
		ok   bool
	}{
		{kind: "StringAttribute", want: "types.String", ok: true},
		{kind: "BoolAttribute", want: "types.Bool", ok: true},
		{kind: "SingleNestedAttribute", want: "types.Object", ok: true},
		{kind: "SetNestedAttribute", want: "types.Set", ok: true},
		{kind: "UnsupportedAttribute", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			t.Parallel()

			got, ok := schemaTypeToModelType(tt.kind)
			if ok != tt.ok {
				t.Fatalf("got ok=%v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
