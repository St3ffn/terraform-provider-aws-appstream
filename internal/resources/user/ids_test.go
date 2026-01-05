// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildID(t *testing.T) {
	require.Equal(t, "AWS_AD|test@example.com", buildID("AWS_AD", "test@example.com"))
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name                   string
		id                     string
		wantAuthenticationType string
		wantUserName           string
		wantErr                bool
	}{
		{
			name:                   "valid_id",
			id:                     "USERPOOL|test@example.com",
			wantAuthenticationType: "USERPOOL",
			wantUserName:           "test@example.com",
			wantErr:                false,
		},
		{
			name:    "missing_separator",
			id:      "USERPOOL-test@example.com",
			wantErr: true,
		},
		{
			name:    "empty_string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "empty_authentication_type",
			id:      "|test@example.com",
			wantErr: true,
		},
		{
			name:    "empty_user_name",
			id:      "USERPOOL|",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authenticationType, userName, err := parseID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if authenticationType != tt.wantAuthenticationType || userName != tt.wantUserName {
				t.Fatalf(
					"parseID(%q) = (%q, %q), want (%q, %q)",
					tt.id, authenticationType, userName, tt.wantAuthenticationType, tt.wantUserName,
				)
			}
		})
	}
}
