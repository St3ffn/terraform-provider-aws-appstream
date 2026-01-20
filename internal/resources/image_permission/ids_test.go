// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_permission

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildID(t *testing.T) {
	require.Equal(t, "test-image|123456789012", buildID("test-image", "123456789012"))
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name                string
		id                  string
		wantName            string
		wantSharedAccountID string
		wantErr             bool
	}{
		{
			name:                "valid_id",
			id:                  "test-image|123456789012",
			wantName:            "test-image",
			wantSharedAccountID: "123456789012",
			wantErr:             false,
		},
		{
			name:    "missing_separator",
			id:      "test-image-123456789012",
			wantErr: true,
		},
		{
			name:    "empty_string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "empty_authentication_type",
			id:      "|123456789012",
			wantErr: true,
		},
		{
			name:    "empty_user_name",
			id:      "test-image|",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, sharedAccountID, err := parseID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if name != tt.wantName || sharedAccountID != tt.wantSharedAccountID {
				t.Fatalf(
					"parseID(%q) = (%q, %q), want (%q, %q)",
					tt.id, name, sharedAccountID, tt.wantName, tt.wantSharedAccountID,
				)
			}
		})
	}
}
