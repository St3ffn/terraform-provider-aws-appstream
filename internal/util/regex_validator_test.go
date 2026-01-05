// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package util

import "testing"

func TestValidateRegexPattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		wantErr bool
	}{
		{
			name:    "valid_simple_regex",
			pattern: "^abc$",
			wantErr: false,
		},
		{
			name:    "valid_character_class",
			pattern: "^[a-zA-Z0-9_-]+$",
			wantErr: false,
		},
		{
			name:    "valid_posix_regex",
			pattern: "^[[:alpha:]]+$",
			wantErr: false,
		},
		{
			name:    "invalid_unclosed_group",
			pattern: "(",
			wantErr: true,
		},
		{
			name:    "invalid_unbalanced_brackets",
			pattern: "[a-z",
			wantErr: true,
		},
		{
			name:    "invalid_escape",
			pattern: "\\",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRegexPattern(tt.pattern)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
