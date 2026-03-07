// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package copied_image

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildID(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"example-image|us-east-1|source-image",
		buildID("example-image", "us-east-1", "source-image"),
	)
}

func TestParseID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                     string
		id                       string
		wantDestinationImageName string
		wantDestinationRegion    string
		wantSourceImageName      string
		wantErr                  bool
	}{
		{
			name:                     "valid_id",
			id:                       "example-image|us-east-1|source-image",
			wantDestinationImageName: "example-image",
			wantDestinationRegion:    "us-east-1",
			wantSourceImageName:      "source-image",
			wantErr:                  false,
		},
		{
			name:    "missing_separator",
			id:      "example-image-us-east-1",
			wantErr: true,
		},
		{
			name:    "empty_string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "empty_name",
			id:      "|us-east-1|source-image",
			wantErr: true,
		},
		{
			name:    "empty_destination_region",
			id:      "example-image||source-image",
			wantErr: true,
		},
		{
			name:    "empty_source_image_name",
			id:      "example-image|us-east-1|",
			wantErr: true,
		},
		{
			name:    "too_many_parts",
			id:      "example-image|us-east-1|source-image|extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			name, destinationRegion, sourceImageName, err := parseID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if name != tt.wantDestinationImageName ||
				destinationRegion != tt.wantDestinationRegion ||
				sourceImageName != tt.wantSourceImageName {
				t.Fatalf(
					"parseID(%q) = (%q, %q, %q), want (%q, %q, %q)",
					tt.id,
					name,
					destinationRegion,
					sourceImageName,
					tt.wantDestinationImageName,
					tt.wantDestinationRegion,
					tt.wantSourceImageName,
				)
			}
		})
	}
}
