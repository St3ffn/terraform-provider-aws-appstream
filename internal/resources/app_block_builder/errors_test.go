// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsUnexpectedAppBlockBuilderStateError(t *testing.T) {
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
			err:  errUnexpectedAppBlockBuilderState,
			want: true,
		},
		{
			name: "wrapped_error",
			err:  fmt.Errorf("wrapped: %w", errUnexpectedAppBlockBuilderState),
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

			got := isUnexpectedAppBlockBuilderStateError(tt.err)
			require.Equal(t, tt.want, got)
		})
	}
}
