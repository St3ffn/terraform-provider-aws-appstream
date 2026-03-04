// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package gomod

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

const (
	// AppStreamModulePath is the Go module path for the AWS AppStream SDK service package.
	AppStreamModulePath = "github.com/aws/aws-sdk-go-v2/service/appstream"
)

// AppStreamVersion reads the AppStream SDK version from the provided go.mod file.
func AppStreamVersion(goModPath string) (string, error) {
	return ModuleVersion(goModPath, AppStreamModulePath)
}

// ModuleVersion reads a module version from the provided go.mod file. If a replace
// directive exists for the module and includes a replacement version, that version
// is returned as the effective version.
func ModuleVersion(goModPath, modulePath string) (string, error) {
	goModPath = strings.TrimSpace(goModPath)
	modulePath = strings.TrimSpace(modulePath)

	if goModPath == "" {
		return "", errors.New("go.mod path must not be empty")
	}
	if modulePath == "" {
		return "", errors.New("module path must not be empty")
	}

	goModPath, err := sanitizeGoModPath(goModPath)
	if err != nil {
		return "", err
	}

	// #nosec G304 -- goModPath is sanitized by sanitizeGoModPath to an absolute path inside repo root with basename go.mod.
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("read go.mod: %w", err)
	}

	f, err := modfile.Parse(goModPath, content, nil)
	if err != nil {
		return "", fmt.Errorf("parse go.mod: %w", err)
	}

	var version string
	for _, req := range f.Require {
		if req.Mod.Path == modulePath {
			version = req.Mod.Version
			break
		}
	}

	if version == "" {
		return "", fmt.Errorf("module %q not found in go.mod requires", modulePath)
	}

	for _, rep := range f.Replace {
		if rep.Old.Path != modulePath {
			continue
		}
		if rep.New.Version != "" {
			version = rep.New.Version
		}
		break
	}

	if !semver.IsValid(version) {
		return "", fmt.Errorf("module %q has invalid semver version %q", modulePath, version)
	}

	return version, nil
}

func sanitizeGoModPath(input string) (string, error) {
	cleanPath := filepath.Clean(input)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("resolve absolute go.mod path: %w", err)
	}

	safeRoot, err := findRepositoryRoot()
	if err != nil {
		return "", err
	}

	safeRoot = filepath.Clean(safeRoot)
	safeRootPrefix := safeRoot + string(filepath.Separator)
	if absPath != safeRoot && !strings.HasPrefix(absPath, safeRootPrefix) {
		return "", fmt.Errorf("unsafe go.mod path %q: must be inside repository root %q", absPath, safeRoot)
	}

	if filepath.Base(absPath) != "go.mod" {
		return "", fmt.Errorf("unsafe go.mod path %q: expected file name go.mod", absPath)
	}

	return absPath, nil
}

func findRepositoryRoot() (string, error) {
	startDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve current working directory: %w", err)
	}

	dir := filepath.Clean(startDir)
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("repository root not found (missing .git)")
		}
		dir = parent
	}
}
