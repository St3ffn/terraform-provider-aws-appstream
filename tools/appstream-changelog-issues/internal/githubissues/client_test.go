// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package githubissues

import "testing"

func TestNewClient_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		owner string
		repo  string
		token *Token
	}{
		{name: "missing owner", repo: "r", token: NewToken("t")},
		{name: "missing repo", owner: "o", token: NewToken("t")},
		{name: "missing token", owner: "o", repo: "r"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewClient(tt.owner, tt.repo, tt.token)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}
