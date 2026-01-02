// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_application_fleet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildID(t *testing.T) {
	require.Equal(t,
		"fleet1|arn:aws:appstream:eu-west-1:123456789012:application/test-app",
		buildID("fleet1", "arn:aws:appstream:eu-west-1:123456789012:application/test-app"),
	)
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name            string
		id              string
		wantFleet       string
		wantApplication string
		wantErr         bool
	}{
		{
			name:            "valid_id",
			id:              "fleet1|arn:aws:appstream:eu-west-1:123456789012:application/test-app",
			wantFleet:       "fleet1",
			wantApplication: "arn:aws:appstream:eu-west-1:123456789012:application/test-app",
			wantErr:         false,
		},
		{
			name:    "missing_separator",
			id:      "fleet1-arn:aws:appstream:eu-west-1:123456789012:application/test-app",
			wantErr: true,
		},
		{
			name:    "empty_string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "empty_fleet_name",
			id:      "|arn:aws:appstream:eu-west-1:123456789012:application/test-app",
			wantErr: true,
		},
		{
			name:    "empty_application_name",
			id:      "fleet1|",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fleet, application, err := parseID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if fleet != tt.wantFleet || application != tt.wantApplication {
				t.Fatalf(
					"parseID(%q) = (%q, %q), want (%q, %q)",
					tt.id, fleet, application, tt.wantFleet, tt.wantApplication,
				)
			}
		})
	}
}
