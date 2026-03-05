// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

type field struct {
	Name     string
	Type     string
	Tag      string
	Comment  string
	Required bool
	Optional bool
	Computed bool
	Remote   *bool
}

type modelType struct {
	Name   string
	Fields []field
}

// Options configures schema-to-model artifact generation.
type Options struct {
	RootPath       string
	SchemaPatterns []string
	Verbose        bool
}
