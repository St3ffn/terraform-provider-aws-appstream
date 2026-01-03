// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildID(t *testing.T) {
	require.Equal(
		t,
		"stack1|USERPOOL|user@example.com",
		buildID("stack1", "USERPOOL", "user@example.com"),
	)
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		wantStack string
		wantAuth  string
		wantUser  string
		wantErr   bool
	}{
		{
			name:      "valid_id",
			id:        "stack1|USERPOOL|user@example.com",
			wantStack: "stack1",
			wantAuth:  "USERPOOL",
			wantUser:  "user@example.com",
			wantErr:   false,
		},
		{
			name:    "missing_separators",
			id:      "stack1-USERPOOL-user@example.com",
			wantErr: true,
		},
		{
			name:    "empty_string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "empty_stack_name",
			id:      "|USERPOOL|user@example.com",
			wantErr: true,
		},
		{
			name:    "empty_authentication_type",
			id:      "stack1||user@example.com",
			wantErr: true,
		},
		{
			name:    "empty_user_name",
			id:      "stack1|USERPOOL|",
			wantErr: true,
		},
		{
			name:    "too_few_parts",
			id:      "stack1|USERPOOL",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack, auth, user, err := parseID(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantStack, stack)
			require.Equal(t, tt.wantAuth, auth)
			require.Equal(t, tt.wantUser, user)
		})
	}
}
