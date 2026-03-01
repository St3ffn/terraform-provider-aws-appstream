// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import (
	"go/ast"
	"go/token"
	"strings"
)

func fieldStringValue(cl *ast.CompositeLit, key string) string {
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		k, ok := kv.Key.(*ast.Ident)
		if !ok || k.Name != key {
			continue
		}

		v, _ := evalStringExpr(kv.Value)
		return v
	}

	return ""
}

func fieldBoolValue(cl *ast.CompositeLit, key string) bool {
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		k, ok := kv.Key.(*ast.Ident)
		if !ok || k.Name != key {
			continue
		}

		v, _ := evalBoolExpr(kv.Value)
		return v
	}

	return false
}

func evalStringExpr(expr ast.Expr) (string, bool) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		return literalString(v)
	case *ast.BinaryExpr:
		if v.Op != token.ADD {
			return "", false
		}

		left, ok := evalStringExpr(v.X)
		if !ok {
			return "", false
		}

		right, ok := evalStringExpr(v.Y)
		if !ok {
			return "", false
		}

		return left + right, true
	case *ast.ParenExpr:
		return evalStringExpr(v.X)
	default:
		return "", false
	}
}

func evalBoolExpr(expr ast.Expr) (bool, bool) {
	switch v := expr.(type) {
	case *ast.Ident:
		switch v.Name {
		case "true":
			return true, true
		case "false":
			return false, true
		default:
			return false, false
		}
	case *ast.ParenExpr:
		return evalBoolExpr(v.X)
	default:
		return false, false
	}
}

func literalStringExpr(expr ast.Expr) (string, bool) {
	bl, ok := expr.(*ast.BasicLit)
	if !ok {
		return "", false
	}
	return literalString(bl)
}

func selectorName(expr ast.Expr) (string, bool) {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		if t.Sel == nil {
			return "", false
		}
		return t.Sel.Name, true
	default:
		return "", false
	}
}

func literalString(expr *ast.BasicLit) (string, bool) {
	if expr == nil || expr.Kind != token.STRING {
		return "", false
	}

	return strings.Trim(expr.Value, "\""), true
}
