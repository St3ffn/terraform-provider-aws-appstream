// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import (
	"fmt"
	"go/ast"
)

func findTopLevelAttributes(file *ast.File) (*ast.CompositeLit, error) {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name != "Schema" || fn.Body == nil {
			continue
		}

		for _, stmt := range fn.Body.List {
			assign, ok := stmt.(*ast.AssignStmt)
			if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
				continue
			}

			lhsSel, ok := assign.Lhs[0].(*ast.SelectorExpr)
			if !ok || lhsSel.Sel == nil || lhsSel.Sel.Name != "Schema" {
				continue
			}

			schemaLit, ok := assign.Rhs[0].(*ast.CompositeLit)
			if !ok {
				continue
			}

			attrsLit := findFieldCompositeLiteral(schemaLit, "Attributes")
			if attrsLit != nil {
				return attrsLit, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find resp.Schema.Attributes map in schema function")
}

func parseFields(
	attrsMap *ast.CompositeLit,
	commentMap ast.CommentMap,
	path []string,
	rootTypeName string,
	nestedModels *[]modelType,
) ([]field, error) {
	fields := make([]field, 0, len(attrsMap.Elts))

	for _, elt := range attrsMap.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		tfName, ok := literalStringExpr(kv.Key)
		if !ok {
			continue
		}

		attrLit, ok := kv.Value.(*ast.CompositeLit)
		if !ok {
			continue
		}

		attrKind, ok := selectorName(attrLit.Type)
		if !ok {
			continue
		}

		description := fieldStringValue(attrLit, "MarkdownDescription")
		if description == "" {
			description = fieldStringValue(attrLit, "Description")
		}
		required := fieldBoolValue(attrLit, "Required")
		optional := fieldBoolValue(attrLit, "Optional")
		computed := fieldBoolValue(attrLit, "Computed")
		remote := remoteOverrideFromCommentMap(commentMap, kv)

		goType, ok := schemaTypeToModelType(attrKind)
		if !ok {
			return nil, fmt.Errorf(
				"unsupported schema attribute kind %q for %q",
				attrKind,
				tfName,
			)
		}

		nestedAttrs := findNestedAttributes(attrLit, attrKind)
		if nestedAttrs != nil {
			nestedPath := append(append([]string{}, path...), tfName)
			nestedName := nestedModelTypeName(rootTypeName, nestedPath)

			nestedFields, err := parseFields(
				nestedAttrs,
				commentMap,
				nestedPath,
				rootTypeName,
				nestedModels,
			)
			if err != nil {
				return nil, err
			}

			*nestedModels = append(*nestedModels, modelType{
				Name:   nestedName,
				Fields: nestedFields,
			})
		}

		fields = append(fields, field{
			Name:     toGoFieldName(tfName),
			Type:     goType,
			Tag:      tfName,
			Comment:  description,
			Required: required,
			Optional: optional,
			Computed: computed,
			Remote:   remote,
		})
	}

	return fields, nil
}

func findNestedAttributes(attrLit *ast.CompositeLit, attrKind string) *ast.CompositeLit {
	switch attrKind {
	case "SingleNestedAttribute":
		return findFieldCompositeLiteral(attrLit, "Attributes")
	case "SetNestedAttribute", "ListNestedAttribute", "MapNestedAttribute":
		nestedObject := findFieldCompositeLiteral(attrLit, "NestedObject")
		if nestedObject == nil {
			return nil
		}

		return findFieldCompositeLiteral(nestedObject, "Attributes")
	default:
		return nil
	}
}

func findFieldCompositeLiteral(cl *ast.CompositeLit, fieldName string) *ast.CompositeLit {
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != fieldName {
			continue
		}

		nested, ok := kv.Value.(*ast.CompositeLit)
		if !ok {
			return nil
		}

		return nested
	}

	return nil
}
