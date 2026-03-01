// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func findSchemaFiles(rootPath string, schemaPatterns []string) ([]string, error) {
	files := make([]string, 0, 32)

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			return nil
		}

		if matchesAnyPattern(filepath.Base(path), schemaPatterns) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk root path: %w", err)
	}

	sort.Strings(files)
	return files, nil
}

func matchesAnyPattern(name string, schemaPatterns []string) bool {
	for _, pattern := range schemaPatterns {
		if strings.ContainsAny(pattern, "*?[]") {
			matched, err := filepath.Match(pattern, name)
			if err == nil && matched {
				return true
			}
			continue
		}

		if name == pattern {
			return true
		}
	}

	return false
}

func generatedTarget(schemaPath string) (outPath, typeName string) {
	base := filepath.Base(schemaPath)
	dir := filepath.Dir(schemaPath)

	switch base {
	case "resource_schema.go":
		return filepath.Join(dir, "resource_model_gen.go"), "resourceModel"
	case "data_source_schema.go":
		return filepath.Join(dir, "data_source_model_gen.go"), "dataSourceModel"
	default:
		prefix := strings.TrimSuffix(base, "_schema.go")
		if prefix == base {
			prefix = strings.TrimSuffix(base, ".go")
		}
		if prefix == "" {
			prefix = "schema"
		}

		typeName = lowerFirst(toGoFieldName(prefix)) + "Model"
		outPath = filepath.Join(dir, prefix+"_model_gen.go")
		return outPath, typeName
	}
}

func generatedDiffTarget(schemaPath string, modelTypeName string) (outPath, diffTypeName, constructorName string) {
	base := filepath.Base(schemaPath)
	dir := filepath.Dir(schemaPath)

	switch base {
	case "resource_schema.go":
		return filepath.Join(dir, "resource_diff_gen.go"), "resourceDiff", "newResourceDiff"
	default:
		prefix := strings.TrimSuffix(base, "_schema.go")
		if prefix == base {
			prefix = strings.TrimSuffix(base, ".go")
		}
		if prefix == "" {
			prefix = "schema"
		}

		diffTypeName = strings.TrimSuffix(modelTypeName, "Model") + "Diff"
		constructorName = "new" + upperFirst(diffTypeName)
		outPath = filepath.Join(dir, prefix+"_diff_gen.go")
		return outPath, diffTypeName, constructorName
	}
}

func shouldGenerateDiff(schemaPath string) bool {
	return filepath.Base(schemaPath) != "data_source_schema.go"
}
