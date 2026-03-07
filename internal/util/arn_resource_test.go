// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
)

func TestARNResourceSuffixOrNil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		value          *string
		service        string
		resourcePrefix string
		want           *string
	}{
		{
			name:           "nil_input_returns_nil",
			value:          nil,
			service:        "appstream",
			resourcePrefix: "image/",
			want:           nil,
		},
		{
			name:           "invalid_arn_returns_nil",
			value:          aws.String("not-an-arn"),
			service:        "appstream",
			resourcePrefix: "image/",
			want:           nil,
		},
		{
			name:           "wrong_service_returns_nil",
			value:          aws.String("arn:aws:iam::123456789012:role/MyRole"),
			service:        "appstream",
			resourcePrefix: "image/",
			want:           nil,
		},
		{
			name:           "wrong_resource_prefix_returns_nil",
			value:          aws.String("arn:aws:appstream:eu-central-1:123456789012:image-builder/example"),
			service:        "appstream",
			resourcePrefix: "image/",
			want:           nil,
		},
		{
			name:           "empty_suffix_returns_nil",
			value:          aws.String("arn:aws:appstream:eu-central-1:123456789012:image/"),
			service:        "appstream",
			resourcePrefix: "image/",
			want:           nil,
		},
		{
			name:           "valid_arn_returns_suffix",
			value:          aws.String("arn:aws:appstream:eu-central-1:123456789012:image/AppStream-WinServer2025-12-18-2025"),
			service:        "appstream",
			resourcePrefix: "image/",
			want:           aws.String("AppStream-WinServer2025-12-18-2025"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ARNResourceSuffixOrNil(tt.value, tt.service, tt.resourcePrefix)

			if tt.want == nil {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			require.Equal(t, *tt.want, *got)
		})
	}
}
