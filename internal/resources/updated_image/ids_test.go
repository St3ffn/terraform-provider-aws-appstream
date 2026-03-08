// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package updated_image

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildID(t *testing.T) {
	t.Parallel()

	require.Equal(t, "existing-image|new-image", buildID("existing-image", "new-image"))
}

func TestParseID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		id                    string
		wantExistingImageName string
		wantNewImageName      string
		wantErr               bool
	}{
		{
			name:                  "valid_id",
			id:                    "existing-image|new-image",
			wantExistingImageName: "existing-image",
			wantNewImageName:      "new-image",
			wantErr:               false,
		},
		{
			name:    "missing_separator",
			id:      "existing-image-new-image",
			wantErr: true,
		},
		{
			name:    "empty_string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "empty_existing_image_name",
			id:      "|new-image",
			wantErr: true,
		},
		{
			name:    "empty_new_image_name",
			id:      "existing-image|",
			wantErr: true,
		},
		{
			name:    "too_many_parts",
			id:      "existing-image|new-image|extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			existingImageName, newImageName, err := parseID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if existingImageName != tt.wantExistingImageName ||
				newImageName != tt.wantNewImageName {
				t.Fatalf(
					"parseID(%q) = (%q, %q), want (%q, %q)",
					tt.id,
					existingImageName,
					newImageName,
					tt.wantExistingImageName,
					tt.wantNewImageName,
				)
			}
		})
	}
}
