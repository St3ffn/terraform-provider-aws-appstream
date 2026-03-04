// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_RequiresOptions(t *testing.T) {
	t.Parallel()

	err := Run(Options{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestRun_GeneratesModelAndResourceDiffFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	schemaPath := filepath.Join(root, "resource_schema.go")

	schema := `package sample

import (
	"context"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resource struct{}

func (r *resource) Schema(_ context.Context, _ tfresource.SchemaRequest, resp *tfresource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the fleet.",
				Required: true,
			},
			"vpc_config": schema.SingleNestedAttribute{
				Description: "VPC configuration.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"subnet_ids": schema.SetAttribute{
						Description: "Subnet IDs.",
						Optional: true,
						ElementType: types.StringType,
					},
				},
			},
			// codegen:has_remote_changes=false
			"desired_state": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"update_behavior": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}
`

	if err := os.WriteFile(schemaPath, []byte(schema), 0o600); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	if err := Run(Options{
		RootPath:       root,
		SchemaPatterns: []string{"resource_schema.go"},
	}); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	outPath := filepath.Join(root, "resource_model_gen.go")
	// #nosec G304 -- fixed generated filename under t.TempDir(), not user-controlled path input.
	out, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}

	content := string(out)
	mustContain := []string{
		"type resourceModel struct",
		"Name of the fleet.",
		"Name types.String `tfsdk:\"name\"`",
		"type resourceModelVPCConfig struct",
	}

	for _, needle := range mustContain {
		if !strings.Contains(content, needle) {
			t.Fatalf("generated file does not contain %q\ncontent:\n%s", needle, content)
		}
	}

	diffPath := filepath.Join(root, "resource_diff_gen.go")
	// #nosec G304 -- fixed generated filename under t.TempDir(), not user-controlled path input.
	diffOut, err := os.ReadFile(diffPath)
	if err != nil {
		t.Fatalf("read generated diff file: %v", err)
	}

	diffContent := string(diffOut)
	diffMustContain := []string{
		"type resourceDiff struct",
		"type changeKind uint8",
		"changeNone changeKind = iota",
		"func (k changeKind) IsUpdated() bool",
		"func (k changeKind) IsCleared() bool",
		"func newResourceDiff(state resourceModel, plan resourceModel) resourceDiff",
		"VPCConfig",
		"func classifyChange(state attr.Value, plan attr.Value) changeKind",
		"func (d resourceDiff) HasRemoteChanges() bool",
		"func (d resourceDiff) RequiresTagApply() bool",
	}

	for _, needle := range diffMustContain {
		if !strings.Contains(diffContent, needle) {
			t.Fatalf("generated diff file does not contain %q\ncontent:\n%s", needle, diffContent)
		}
	}

	if strings.Contains(diffContent, "d.DesiredState.IsChanged()") {
		t.Fatalf("desired_state should be excluded from HasRemoteChanges when codegen:has_remote_changes=false is set")
	}

	if !strings.Contains(diffContent, "d.UpdateBehavior.IsChanged()") {
		t.Fatalf("update_behavior should be included in HasRemoteChanges by default")
	}
}

func TestRun_DataSourceDoesNotGenerateDiffFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	schemaPath := filepath.Join(root, "data_source_schema.go")

	schema := `package sample

import (
	"context"

	tfdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type dataSource struct{}

func (d *dataSource) Schema(_ context.Context, _ tfdatasource.SchemaRequest, resp *tfdatasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"arn": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}
`

	if err := os.WriteFile(schemaPath, []byte(schema), 0o600); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	if err := Run(Options{
		RootPath:       root,
		SchemaPatterns: []string{"data_source_schema.go"},
	}); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	modelPath := filepath.Join(root, "data_source_model_gen.go")
	if _, err := os.Stat(modelPath); err != nil {
		t.Fatalf("expected model file to exist: %v", err)
	}

	diffPath := filepath.Join(root, "data_source_diff_gen.go")
	if _, err := os.Stat(diffPath); err == nil {
		t.Fatalf("did not expect diff file for data source")
	}
}
