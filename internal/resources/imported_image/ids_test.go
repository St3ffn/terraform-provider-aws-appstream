// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package imported_image

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildID(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"example-image|arn:aws:iam::123456789012:role/example|ami-0abc1234def567890",
		buildID(
			"example-image",
			"arn:aws:iam::123456789012:role/example",
			"ami-0abc1234def567890",
		),
	)
}

func TestParseID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		id             string
		wantName       string
		wantIAMRoleARN string
		wantSourceAmi  string
		wantErr        bool
	}{
		{
			name:           "valid_id",
			id:             "example-image|arn:aws:iam::123456789012:role/example|ami-0abc1234def567890",
			wantName:       "example-image",
			wantIAMRoleARN: "arn:aws:iam::123456789012:role/example",
			wantSourceAmi:  "ami-0abc1234def567890",
		},
		{
			name:    "missing_separator",
			id:      "example-image-arn:aws:iam::123456789012:role/example-ami-0abc1234def567890",
			wantErr: true,
		},
		{
			name:    "empty_string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "empty_name",
			id:      "|arn:aws:iam::123456789012:role/example|ami-0abc1234def567890",
			wantErr: true,
		},
		{
			name:    "empty_iam_role_arn",
			id:      "example-image||ami-0abc1234def567890",
			wantErr: true,
		},
		{
			name:    "empty_source_ami_id",
			id:      "example-image|arn:aws:iam::123456789012:role/example|",
			wantErr: true,
		},
		{
			name:    "too_many_parts",
			id:      "example-image|arn:aws:iam::123456789012:role/example|ami-0abc1234def567890|extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			name, iamRoleARN, sourceAmiID, err := parseID(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantName, name)
			require.Equal(t, tt.wantIAMRoleARN, iamRoleARN)
			require.Equal(t, tt.wantSourceAmi, sourceAmiID)
		})
	}
}
