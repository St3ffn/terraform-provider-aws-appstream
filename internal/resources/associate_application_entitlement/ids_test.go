// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_application_entitlement

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildID(t *testing.T) {
	require.Equal(t,
		"stack1|entitlement1|app1",
		buildID("stack1", "entitlement1", "app1"),
	)
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name            string
		id              string
		wantStack       string
		wantEntitlement string
		wantAppID       string
		wantErr         bool
	}{
		{
			name:            "valid_id",
			id:              "stack1|entitlement1|app1",
			wantStack:       "stack1",
			wantEntitlement: "entitlement1",
			wantAppID:       "app1",
			wantErr:         false,
		},
		{
			name:    "missing_separator",
			id:      "stack1-entitlement1-app1",
			wantErr: true,
		},
		{
			name:    "empty_string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "empty_stack_name",
			id:      "|entitlement1|app1",
			wantErr: true,
		},
		{
			name:    "empty_entitlement_name",
			id:      "stack1||app1",
			wantErr: true,
		},
		{
			name:    "empty_application_name",
			id:      "stack1|entitlement1|",
			wantErr: true,
		},
		{
			name:    "empty_stack_name_and_app_name",
			id:      "|entitlement1|",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack, entitlement, appID, err := parseID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if stack != tt.wantStack || entitlement != tt.wantEntitlement || appID != tt.wantAppID {
				t.Fatalf(
					"parseID(%q) = (%q, %q, %q), want (%q, %q, %q)",
					tt.id, stack, entitlement, appID, tt.wantStack, tt.wantEntitlement, tt.wantAppID,
				)
			}
		})
	}
}
