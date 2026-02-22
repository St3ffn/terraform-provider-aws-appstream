// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package githubissues

import (
	"fmt"
	"testing"
)

func TestNewToken(t *testing.T) {
	t.Parallel()

	s := NewToken("  token-value  ")

	if !s.IsSet() {
		t.Fatal("expected token to be set")
	}
	if got := s.value(); got != "token-value" {
		t.Fatalf("expected trimmed value %q, got %q", "token-value", got)
	}
	if got := s.String(); got != "******" {
		t.Fatalf("expected masked string %q, got %q", "******", got)
	}
	if got := s.GoString(); got != "******" {
		t.Fatalf("expected masked gostring %q, got %q", "******", got)
	}
}

func TestNewTokenWhitespaceOnly(t *testing.T) {
	t.Parallel()

	s := NewToken("   \n\t  ")

	if s.IsSet() {
		t.Fatal("expected token to be unset")
	}
	if got := s.value(); got != "" {
		t.Fatalf("expected empty value, got %q", got)
	}
	if got := s.String(); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
	if got := s.GoString(); got != "" {
		t.Fatalf("expected empty gostring, got %q", got)
	}
}

func TestTokenNilReceiver(t *testing.T) {
	t.Parallel()

	var s *Token

	if s.IsSet() {
		t.Fatal("expected nil token to be unset")
	}
	if got := s.value(); got != "" {
		t.Fatalf("expected empty value for nil token, got %q", got)
	}
	if got := s.String(); got != "" {
		t.Fatalf("expected empty string for nil token, got %q", got)
	}
	if got := s.GoString(); got != "" {
		t.Fatalf("expected empty gostring for nil token, got %q", got)
	}
}

func TestTokenFmtMasking(t *testing.T) {
	t.Parallel()

	s := NewToken("token")

	if got := fmt.Sprintf("%v", s); got != "******" {
		t.Fatalf("expected %%v to be masked, got %q", got)
	}
	if got := fmt.Sprintf("%#v", s); got != "******" {
		t.Fatalf("expected %%#v to be masked, got %q", got)
	}
}
