// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
)

func TestImageNameFromImageARN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arn  *string
		want *string
	}{
		{
			name: "valid image arn",
			arn:  aws.String("arn:aws:appstream:eu-central-1:123456789012:image/AppStream-WinServer2025-12-18-2025"),
			want: aws.String("AppStream-WinServer2025-12-18-2025"),
		},
		{
			name: "nil arn",
			arn:  nil,
			want: nil,
		},
		{
			name: "invalid arn resource",
			arn:  aws.String("arn:aws:appstream:eu-central-1:123456789012:image-builder/example"),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := imageNameFromImageARN(tt.arn)

			if tt.want == nil {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			require.Equal(t, *tt.want, *got)
		})
	}
}
