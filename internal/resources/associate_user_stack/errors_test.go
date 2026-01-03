// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	"github.com/stretchr/testify/require"
)

func TestNewUserStackAssociationNotReadyError(t *testing.T) {
	tests := []struct {
		name       string
		input      awstypes.UserStackAssociationError
		wantNil    bool
		wantMsgSub string
	}{
		{
			name: "with_error_message",
			input: awstypes.UserStackAssociationError{
				ErrorCode:    awstypes.UserStackAssociationErrorCodeStackNotFound,
				ErrorMessage: aws.String("stack does not exist"),
			},
			wantNil:    false,
			wantMsgSub: "stack does not exist",
		},
		{
			name: "nil_error_message",
			input: awstypes.UserStackAssociationError{
				ErrorCode: awstypes.UserStackAssociationErrorCodeInternalError,
			},
			wantNil:    false,
			wantMsgSub: "unknown appstream user-stack association error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newUserStackAssociationNotReadyError(tt.input)

			if tt.wantNil {
				require.Nil(t, err)
				return
			}

			require.NotNil(t, err)
			require.True(t, isUserStackAssociationNotReadyError(err))
			require.Contains(t, err.Error(), tt.wantMsgSub)
		})
	}
}

func TestIsUserStackAssociationNotReadyError(t *testing.T) {
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
			err:  errUserStackAssociationNotReady,
			want: true,
		},
		{
			name: "wrapped_error",
			err:  fmt.Errorf("wrapped: %w", errUserStackAssociationNotReady),
			want: true,
		},
		{
			name: "different_error",
			err:  fmt.Errorf("some other error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isUserStackAssociationNotReadyError(tt.err)
			require.Equal(t, tt.want, got)
		})
	}
}
