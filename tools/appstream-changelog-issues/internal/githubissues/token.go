// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package githubissues

import "strings"

// Token stores a GitHub token value and ensures redacted string formatting.
type Token struct {
	v string
}

// NewToken wraps a raw token string in a Token value.
func NewToken(token string) *Token {
	return &Token{v: strings.TrimSpace(token)}
}

// IsSet reports whether a non-empty token value is configured.
func (s *Token) IsSet() bool {
	return s != nil && s.v != ""
}

// String returns the string representation of the current value.
func (s *Token) String() string {
	if s == nil || s.v == "" {
		return ""
	}
	return "******"
}

// GoString returns a redacted debug representation of the token.
func (s *Token) GoString() string {
	return s.String()
}

func (s *Token) value() string {
	if s == nil {
		return ""
	}
	return s.v
}
