// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package report

import (
	"strings"
	"testing"

	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/appstream"
)

func TestBuildIssueContent(t *testing.T) {
	t.Parallel()

	releases := []appstream.Release{
		{
			Version: "v1.54.0",
			Date:    "2026-02-20",
			Notes: []string{
				"**Feature**: Adding new attribute to disable IMDS v1 APIs for fleet, Image Builder and AppBlockBuilder instances.",
			},
		},
		{
			Version: "v1.53.0",
			Date:    "2025-12-18",
			Notes: []string{
				"**Dependency Update**: Updated to the latest SDK module versions",
				"**Feature**: Added support for new operating systems",
			},
		},
	}

	got, err := BuildIssueContent("v1.52.4", releases, "https://example.invalid/changelog")
	if err != nil {
		t.Fatalf("BuildIssueContent returned error: %v", err)
	}

	if got.Title != IssueTitle {
		t.Fatalf("expected title %q, got %q", IssueTitle, got.Title)
	}

	wantContains := []string{
		IssueMarker,
		"Current `github.com/aws/aws-sdk-go-v2/service/appstream` version: `v1.52.4`",
		"Latest feature release detected: `v1.54.0`",
		"Feature releases found: `2`",
		"### v1.54.0 (2026-02-20)",
		"### v1.53.0 (2025-12-18)",
		"- **Feature**: Added support for new operating systems",
		"## Source",
		"- https://example.invalid/changelog",
	}

	for _, want := range wantContains {
		if !strings.Contains(got.Body, want) {
			t.Fatalf("expected body to contain %q, got:\n%s", want, got.Body)
		}
	}
}

func TestBuildIssueContent_DefaultSourceURL(t *testing.T) {
	t.Parallel()

	releases := []appstream.Release{
		{
			Version: "v1.54.0",
			Date:    "2026-02-20",
			Notes:   []string{"**Feature**: x"},
		},
	}

	got, err := BuildIssueContent("v1.53.0", releases, "")
	if err != nil {
		t.Fatalf("BuildIssueContent returned error: %v", err)
	}

	if !strings.Contains(got.Body, DefaultSourceURL) {
		t.Fatalf("expected body to contain default source URL %q, got:\n%s", DefaultSourceURL, got.Body)
	}
}

func TestBuildIssueContent_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		currentVersion string
		releases       []appstream.Release
	}{
		{
			name:           "invalid current version",
			currentVersion: "1.54.0",
			releases: []appstream.Release{
				{Version: "v1.54.0", Date: "2026-02-20", Notes: []string{"**Feature**: x"}},
			},
		},
		{
			name:           "empty releases",
			currentVersion: "v1.54.0",
			releases:       nil,
		},
		{
			name:           "invalid release version",
			currentVersion: "v1.54.0",
			releases: []appstream.Release{
				{Version: "1.54.1", Date: "2026-02-20", Notes: []string{"**Feature**: x"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := BuildIssueContent(tt.currentVersion, tt.releases, "")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestBuildIssueContentWithOptions_CustomValues(t *testing.T) {
	t.Parallel()

	releases := []appstream.Release{
		{
			Version: "v1.54.0",
			Date:    "2026-02-20",
			Notes:   []string{"**Feature**: x"},
		},
	}

	got, err := BuildIssueContentWithOptions("v1.53.0", releases, BuildOptions{
		IssueTitle:  "Custom title",
		IssueMarker: "<!-- custom-marker -->",
		SourceURL:   "https://example.invalid/custom",
	})
	if err != nil {
		t.Fatalf("BuildIssueContentWithOptions returned error: %v", err)
	}

	if got.Title != "Custom title" {
		t.Fatalf("expected custom title, got %q", got.Title)
	}
	if !strings.Contains(got.Body, "<!-- custom-marker -->") {
		t.Fatalf("expected body to contain custom marker, got:\n%s", got.Body)
	}
	if !strings.Contains(got.Body, "https://example.invalid/custom") {
		t.Fatalf("expected body to contain custom source URL, got:\n%s", got.Body)
	}
}
