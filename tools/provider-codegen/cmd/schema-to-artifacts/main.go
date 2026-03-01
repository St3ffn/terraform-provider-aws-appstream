// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/st3ffn/terraform-provider-aws-appstream/tools/provider-codegen/internal/schematomodel"
)

type stringSliceFlag []string

func (f *stringSliceFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *stringSliceFlag) Set(value string) error {
	v := strings.TrimSpace(value)
	if v == "" {
		return nil
	}

	*f = append(*f, v)
	return nil
}

func main() {
	var (
		rootPath       string
		schemaPatterns stringSliceFlag
		verbose        bool
	)

	flag.StringVar(&rootPath, "root", "", "Root path to scan for schema files")
	flag.Var(&schemaPatterns, "schema-pattern", "Schema filename pattern (repeatable), e.g. resource_schema.go")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	log.SetFlags(0)

	if rootPath == "" || len(schemaPatterns) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "schema-to-artifacts: --root and at least one --schema-pattern are required")
		os.Exit(2)
	}

	if err := schematomodel.Run(schematomodel.Options{
		RootPath:       rootPath,
		SchemaPatterns: schemaPatterns,
		Verbose:        verbose,
	}); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "schema-to-artifacts: %v\n", err)
		os.Exit(1)
	}
}
