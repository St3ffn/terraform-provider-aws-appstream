// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsContextCanceled(t *testing.T) {
	type testCase struct {
		name string
		err  error
		want bool
	}

	tests := []testCase{
		{
			name: "context_canceled/matches_canceled",
			err:  context.Canceled,
			want: true,
		},
		{
			name: "context_canceled/matches_deadline_exceeded",
			err:  context.DeadlineExceeded,
			want: true,
		},
		{
			name: "context_canceled/other_error_no_match",
			err:  errors.New("other"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsContextCanceled(tt.err)
			require.Equal(t, tt.want, got)
		})
	}
}
