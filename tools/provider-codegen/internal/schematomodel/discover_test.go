// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMatchesAnyPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fileName string
		patterns []string
		want     bool
	}{
		{
			name:     "exact_match",
			fileName: "resource_schema.go",
			patterns: []string{"resource_schema.go"},
			want:     true,
		},
		{
			name:     "glob_match",
			fileName: "data_source_schema.go",
			patterns: []string{"*_schema.go"},
			want:     true,
		},
		{
			name:     "invalid_glob_is_ignored",
			fileName: "resource_schema.go",
			patterns: []string{"["},
			want:     false,
		},
		{
			name:     "no_match",
			fileName: "resource_model.go",
			patterns: []string{"*_schema.go"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := matchesAnyPattern(tt.fileName, tt.patterns)
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratedTarget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		schemaPath   string
		wantOutBase  string
		wantTypeName string
	}{
		{
			name:         "resource_schema",
			schemaPath:   filepath.Join("tmp", "resource_schema.go"),
			wantOutBase:  "resource_model_gen.go",
			wantTypeName: "resourceModel",
		},
		{
			name:         "data_source_schema",
			schemaPath:   filepath.Join("tmp", "data_source_schema.go"),
			wantOutBase:  "data_source_model_gen.go",
			wantTypeName: "dataSourceModel",
		},
		{
			name:         "fallback",
			schemaPath:   filepath.Join("tmp", "stack_theme_schema.go"),
			wantOutBase:  "stack_theme_model_gen.go",
			wantTypeName: "stackThemeModel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			outPath, typeName := generatedTarget(tt.schemaPath)
			if filepath.Base(outPath) != tt.wantOutBase {
				t.Fatalf("got out base %q, want %q", filepath.Base(outPath), tt.wantOutBase)
			}
			if typeName != tt.wantTypeName {
				t.Fatalf("got type name %q, want %q", typeName, tt.wantTypeName)
			}
		})
	}
}

func TestGeneratedDiffTarget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		schemaPath      string
		modelTypeName   string
		wantOutBase     string
		wantDiffType    string
		wantConstructor string
	}{
		{
			name:            "resource_schema",
			schemaPath:      filepath.Join("tmp", "resource_schema.go"),
			modelTypeName:   "resourceModel",
			wantOutBase:     "resource_diff_gen.go",
			wantDiffType:    "resourceDiff",
			wantConstructor: "newResourceDiff",
		},
		{
			name:            "fallback",
			schemaPath:      filepath.Join("tmp", "stack_theme_schema.go"),
			modelTypeName:   "stackThemeModel",
			wantOutBase:     "stack_theme_diff_gen.go",
			wantDiffType:    "stackThemeDiff",
			wantConstructor: "newStackThemeDiff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			outPath, diffTypeName, constructorName := generatedDiffTarget(tt.schemaPath, tt.modelTypeName)
			if filepath.Base(outPath) != tt.wantOutBase {
				t.Fatalf("got out base %q, want %q", filepath.Base(outPath), tt.wantOutBase)
			}
			if diffTypeName != tt.wantDiffType {
				t.Fatalf("got diff type %q, want %q", diffTypeName, tt.wantDiffType)
			}
			if constructorName != tt.wantConstructor {
				t.Fatalf("got constructor %q, want %q", constructorName, tt.wantConstructor)
			}
		})
	}
}

func TestShouldGenerateDiff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		schemaPath string
		want       bool
	}{
		{
			name:       "resource_schema_true",
			schemaPath: filepath.Join("tmp", "resource_schema.go"),
			want:       true,
		},
		{
			name:       "data_source_schema_false",
			schemaPath: filepath.Join("tmp", "data_source_schema.go"),
			want:       false,
		},
		{
			name:       "fallback_true",
			schemaPath: filepath.Join("tmp", "custom_schema.go"),
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := shouldGenerateDiff(tt.schemaPath)
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindSchemaFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mustWrite := func(rel string) {
		t.Helper()

		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(path, []byte("package x\n"), 0o600); err != nil {
			t.Fatalf("write: %v", err)
		}
	}

	mustWrite(filepath.Join("a", "resource_schema.go"))
	mustWrite(filepath.Join("b", "data_source_schema.go"))
	mustWrite(filepath.Join("c", "resource_model.go"))

	got, err := findSchemaFiles(root, []string{"resource_schema.go", "data_source_schema.go"})
	if err != nil {
		t.Fatalf("findSchemaFiles returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("got %d files, want 2 (%v)", len(got), got)
	}

	if filepath.Base(got[0]) != "resource_schema.go" || filepath.Base(got[1]) != "data_source_schema.go" {
		t.Fatalf("unexpected files order/content: %v", got)
	}
}
