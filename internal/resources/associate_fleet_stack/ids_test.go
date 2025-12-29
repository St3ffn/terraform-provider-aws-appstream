// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_fleet_stack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildID(t *testing.T) {
	require.Equal(t, "fleet1|stack1", buildID("fleet1", "stack1"))
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		wantFleet string
		wantStack string
		wantErr   bool
	}{
		{
			name:      "valid_id",
			id:        "fleet1|stack1",
			wantFleet: "fleet1",
			wantStack: "stack1",
			wantErr:   false,
		},
		{
			name:    "missing_separator",
			id:      "fleet1-stack1",
			wantErr: true,
		},
		{
			name:    "empty_string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "empty_fleet_name",
			id:      "|stack1",
			wantErr: true,
		},
		{
			name:    "empty_stack_name",
			id:      "fleet1|",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fleet, stack, err := parseID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if fleet != tt.wantFleet || stack != tt.wantStack {
				t.Fatalf(
					"parseID(%q) = (%q, %q), want (%q, %q)",
					tt.id, fleet, stack, tt.wantFleet, tt.wantStack,
				)
			}
		})
	}
}
