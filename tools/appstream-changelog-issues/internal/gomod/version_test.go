// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package gomod

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppStreamVersion(t *testing.T) {
	t.Parallel()

	goModPath := writeTempGoMod(t, `
module example.com/test

go 1.26

require (
	github.com/aws/aws-sdk-go-v2/service/appstream v1.54.0
)
`)

	got, err := AppStreamVersion(goModPath)
	if err != nil {
		t.Fatalf("AppStreamVersion returned error: %v", err)
	}

	if got != "v1.54.0" {
		t.Fatalf("expected v1.54.0, got %q", got)
	}
}

func TestModuleVersion_WithReplaceVersion(t *testing.T) {
	t.Parallel()

	goModPath := writeTempGoMod(t, `
module example.com/test

go 1.26

require (
	github.com/aws/aws-sdk-go-v2/service/appstream v1.54.0
)

replace github.com/aws/aws-sdk-go-v2/service/appstream => github.com/aws/aws-sdk-go-v2/service/appstream v1.55.0
`)

	got, err := ModuleVersion(goModPath, AppStreamModulePath)
	if err != nil {
		t.Fatalf("ModuleVersion returned error: %v", err)
	}

	if got != "v1.55.0" {
		t.Fatalf("expected replacement version v1.55.0, got %q", got)
	}
}

func TestModuleVersion_WithLocalReplaceKeepsRequiredVersion(t *testing.T) {
	t.Parallel()

	goModPath := writeTempGoMod(t, `
module example.com/test

go 1.26

require github.com/aws/aws-sdk-go-v2/service/appstream v1.54.0

replace github.com/aws/aws-sdk-go-v2/service/appstream => ../appstream-local
`)

	got, err := ModuleVersion(goModPath, AppStreamModulePath)
	if err != nil {
		t.Fatalf("ModuleVersion returned error: %v", err)
	}

	if got != "v1.54.0" {
		t.Fatalf("expected required version v1.54.0, got %q", got)
	}
}

func TestModuleVersion_MissingModule(t *testing.T) {
	t.Parallel()

	goModPath := writeTempGoMod(t, `
module example.com/test

go 1.26

require github.com/aws/aws-sdk-go-v2 v1.39.0
`)

	_, err := ModuleVersion(goModPath, AppStreamModulePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestModuleVersion_InvalidGoMod(t *testing.T) {
	t.Parallel()

	goModPath := writeTempGoMod(t, "not-a-valid-go-mod")

	_, err := ModuleVersion(goModPath, AppStreamModulePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "parse go.mod") {
		t.Fatalf("expected parse error, got %v", err)
	}
}

func TestModuleVersion_PathOutsideRepositoryRoot(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(path, []byte("module outside/test\n"), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	_, err := ModuleVersion(path, AppStreamModulePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unsafe go.mod path") {
		t.Fatalf("expected unsafe path error, got %v", err)
	}
}

func writeTempGoMod(t *testing.T, content string) string {
	t.Helper()

	repoRoot, err := findRepositoryRoot()
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}

	tmpDir, err := os.MkdirTemp(repoRoot, "gomod-test-")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	path := filepath.Join(tmpDir, "go.mod")

	err = os.WriteFile(path, []byte(strings.TrimSpace(content)+"\n"), 0o600)
	if err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	return path
}
