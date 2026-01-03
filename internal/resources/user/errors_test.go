// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsUserNotYetVisibleError(t *testing.T) {
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
			err:  errUserNotYetVisible,
			want: true,
		},
		{
			name: "wrapped_error",
			err:  fmt.Errorf("wrapped: %w", errUserNotYetVisible),
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
			got := isUserNotYetVisibleError(tt.err)
			require.Equal(t, tt.want, got)
		})
	}
}
