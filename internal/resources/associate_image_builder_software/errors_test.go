// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsUnexpectedImageBuilderStateError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil_error",
			err:  nil,
			want: false,
		},
		{
			name: "exact_error",
			err:  errUnexpectedImageBuilderState,
			want: true,
		},
		{
			name: "wrapped_error",
			err:  fmt.Errorf("wrapped: %w", errUnexpectedImageBuilderState),
			want: true,
		},
		{
			name: "different_error",
			err:  errors.New("some other error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isUnexpectedImageBuilderStateError(tt.err)
			require.Equal(t, tt.want, got)
		})
	}
}
