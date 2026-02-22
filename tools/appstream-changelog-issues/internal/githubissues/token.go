// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package githubissues

import "strings"

type Token struct {
	v string
}

func NewToken(token string) *Token {
	return &Token{v: strings.TrimSpace(token)}
}

func (s *Token) IsSet() bool {
	return s != nil && s.v != ""
}

func (s *Token) String() string {
	if s == nil || s.v == "" {
		return ""
	}
	return "******"
}

func (s *Token) GoString() string {
	return s.String()
}

func (s *Token) value() string {
	if s == nil {
		return ""
	}
	return s.v
}
