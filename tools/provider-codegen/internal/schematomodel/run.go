// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Run(opts Options) error {
	if strings.TrimSpace(opts.RootPath) == "" || len(opts.SchemaPatterns) == 0 {
		return fmt.Errorf("root path and at least one schema pattern are required")
	}

	schemaFiles, err := findSchemaFiles(opts.RootPath, opts.SchemaPatterns)
	if err != nil {
		return err
	}

	if len(schemaFiles) == 0 {
		return fmt.Errorf(
			"no schema files found in %q for patterns %q",
			opts.RootPath,
			strings.Join(opts.SchemaPatterns, ","),
		)
	}

	log.Printf("matched schema files: %d", len(schemaFiles))
	if opts.Verbose {
		for _, schemaPath := range schemaFiles {
			log.Printf("matched schema: %s", schemaPath)
		}
	}

	modelCount := 0
	diffCount := 0

	for _, schemaPath := range schemaFiles {
		modelPath, modelTypeName := generatedTarget(schemaPath)
		diffEnabled := shouldGenerateDiff(schemaPath)
		diffPath, diffTypeName, constructorName := "", "", ""
		if diffEnabled {
			diffPath, diffTypeName, constructorName = generatedDiffTarget(schemaPath, modelTypeName)
		}
		if err = generateArtifacts(
			schemaPath,
			modelPath,
			modelTypeName,
			diffEnabled,
			diffPath,
			diffTypeName,
			constructorName,
			opts.Verbose,
		); err != nil {
			return fmt.Errorf("%s: %w", schemaPath, err)
		}
		modelCount++
		if diffEnabled {
			diffCount++
		}
	}

	log.Printf("generated model files: %d", modelCount)
	log.Printf("generated diff files: %d", diffCount)

	return nil
}

func generateArtifacts(
	schemaPath string,
	modelPath string,
	modelTypeName string,
	diffEnabled bool,
	diffPath string,
	diffTypeName string,
	constructorName string,
	verbose bool,
) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, schemaPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse schema: %w", err)
	}
	commentMap := ast.NewCommentMap(fset, file, file.Comments)

	attrs, err := findTopLevelAttributes(file)
	if err != nil {
		return err
	}

	rootModel := modelType{
		Name: modelTypeName,
	}

	nestedModels := make([]modelType, 0, 8)
	fields, err := parseFields(attrs, commentMap, nil, modelTypeName, &nestedModels)
	if err != nil {
		return err
	}
	rootModel.Fields = fields

	modelOut, err := render(file.Name.Name, rootModel, nestedModels)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(modelPath), 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	if err := os.WriteFile(modelPath, modelOut, 0o600); err != nil {
		return fmt.Errorf("write model output: %w", err)
	}
	if verbose {
		log.Printf("generated file: %s", modelPath)
	}

	if !diffEnabled {
		return nil
	}

	diffOut, err := renderDiff(file.Name.Name, rootModel, diffTypeName, constructorName)
	if err != nil {
		return err
	}

	if err := os.WriteFile(diffPath, diffOut, 0o600); err != nil {
		return fmt.Errorf("write diff output: %w", err)
	}
	if verbose {
		log.Printf("generated file: %s", diffPath)
	}

	return nil
}
