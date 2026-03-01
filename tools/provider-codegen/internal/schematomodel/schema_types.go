// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

func schemaTypeToModelType(attrKind string) (string, bool) {
	switch attrKind {
	case "StringAttribute":
		return "types.String", true
	case "BoolAttribute":
		return "types.Bool", true
	case "Int32Attribute":
		return "types.Int32", true
	case "Int64Attribute":
		return "types.Int64", true
	case "Float32Attribute":
		return "types.Float32", true
	case "Float64Attribute":
		return "types.Float64", true
	case "NumberAttribute":
		return "types.Number", true
	case "MapAttribute":
		return "types.Map", true
	case "SetAttribute", "SetNestedAttribute":
		return "types.Set", true
	case "ListAttribute", "ListNestedAttribute":
		return "types.List", true
	case "ObjectAttribute", "SingleNestedAttribute":
		return "types.Object", true
	case "MapNestedAttribute":
		return "types.Map", true
	case "DynamicAttribute":
		return "types.Dynamic", true
	default:
		return "", false
	}
}
