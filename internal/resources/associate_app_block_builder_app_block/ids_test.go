// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_app_block_builder_app_block

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildID(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"builder1|arn:aws:appstream:eu-central-1:123456789012:app-block/example",
		buildID(
			"builder1",
			"arn:aws:appstream:eu-central-1:123456789012:app-block/example",
		),
	)
}

func TestParseID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		id                  string
		wantAppBlockBuilder string
		wantAppBlockARN     string
		wantErr             bool
	}{
		{
			name:                "valid_id",
			id:                  "builder1|arn:aws:appstream:eu-central-1:123456789012:app-block/example",
			wantAppBlockBuilder: "builder1",
			wantAppBlockARN:     "arn:aws:appstream:eu-central-1:123456789012:app-block/example",
			wantErr:             false,
		},
		{
			name:    "missing_separator",
			id:      "builder1-arn",
			wantErr: true,
		},
		{
			name:    "empty_string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "empty_app_block_builder_name",
			id:      "|arn:aws:appstream:eu-central-1:123456789012:app-block/example",
			wantErr: true,
		},
		{
			name:    "empty_app_block_arn",
			id:      "builder1|",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builderName, appBlockARN, err := parseID(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantAppBlockBuilder, builderName)
			require.Equal(t, tt.wantAppBlockARN, appBlockARN)
		})
	}
}
