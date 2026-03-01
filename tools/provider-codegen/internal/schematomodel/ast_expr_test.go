// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import (
	"go/ast"
	"go/parser"
	"testing"
)

func TestEvalStringExpr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		expr string
		want string
		ok   bool
	}{
		{
			name: "basic_literal",
			expr: `"hello"`,
			want: "hello",
			ok:   true,
		},
		{
			name: "concatenation",
			expr: `"hello" + " " + "world"`,
			want: "hello world",
			ok:   true,
		},
		{
			name: "parenthesized_concatenation",
			expr: `("hello" + " ") + "world"`,
			want: "hello world",
			ok:   true,
		},
		{
			name: "not_a_string_expr",
			expr: `1 + 2`,
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parsed, err := parser.ParseExpr(tt.expr)
			if err != nil {
				t.Fatalf("parse expr: %v", err)
			}

			got, ok := evalStringExpr(parsed)
			if ok != tt.ok {
				t.Fatalf("got ok=%v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEvalBoolExpr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		expr ast.Expr
		want bool
		ok   bool
	}{
		{
			name: "true_ident",
			expr: &ast.Ident{Name: "true"},
			want: true,
			ok:   true,
		},
		{
			name: "false_ident",
			expr: &ast.Ident{Name: "false"},
			want: false,
			ok:   true,
		},
		{
			name: "unsupported_ident",
			expr: &ast.Ident{Name: "x"},
			ok:   false,
		},
		{
			name: "paren_true",
			expr: &ast.ParenExpr{X: &ast.Ident{Name: "true"}},
			want: true,
			ok:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := evalBoolExpr(tt.expr)
			if ok != tt.ok {
				t.Fatalf("got ok=%v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
