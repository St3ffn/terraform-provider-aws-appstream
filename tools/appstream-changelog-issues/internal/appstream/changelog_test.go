// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package appstream

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestParseChangelog(t *testing.T) {
	t.Parallel()

	markdown := strings.Join([]string{
		"# v1.53.0 (2025-12-18)",
		"",
		"* **Feature**: Added support for something",
		"* **Dependency Update**: Updated SDK modules",
		"",
		"# v1.52.5 (2025-12-08)",
		"",
		"* **Dependency Update**: Updated modules",
		"  with continuation details",
		"",
	}, "\n")

	got, err := ParseChangelog([]byte(markdown))
	if err != nil {
		t.Fatalf("ParseChangelog returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 releases, got %d", len(got))
	}

	if got[0].Version != "v1.53.0" {
		t.Fatalf("expected first version v1.53.0, got %q", got[0].Version)
	}
	if got[0].Date != "2025-12-18" {
		t.Fatalf("expected first date 2025-12-18, got %q", got[0].Date)
	}
	if len(got[0].Notes) != 2 {
		t.Fatalf("expected 2 notes in first release, got %d", len(got[0].Notes))
	}

	if got[1].Version != "v1.52.5" {
		t.Fatalf("expected second version v1.52.5, got %q", got[1].Version)
	}
	if len(got[1].Notes) != 1 {
		t.Fatalf("expected 1 note in second release, got %d", len(got[1].Notes))
	}
	if !strings.Contains(got[1].Notes[0], "continuation details") {
		t.Fatalf("expected wrapped note to include continuation details, got %q", got[1].Notes[0])
	}
}

func TestParseChangelog_NoReleaseSections(t *testing.T) {
	t.Parallel()

	_, err := ParseChangelog([]byte("no headings here"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewerThan(t *testing.T) {
	t.Parallel()

	releases := []Release{
		{Version: "v1.53.1"},
		{Version: "v1.53.0"},
		{Version: "v1.52.5"},
	}

	got, err := NewerThan(releases, "v1.53.0")
	if err != nil {
		t.Fatalf("NewerThan returned error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 release newer than current, got %d", len(got))
	}
	if got[0].Version != "v1.53.1" {
		t.Fatalf("expected v1.53.1, got %q", got[0].Version)
	}
}

func TestFetchChangelog(t *testing.T) {
	t.Parallel()

	const body = "# v1.0.0 (2025-01-01)\n\n* note\n"

	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method %s", r.Method)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	got, err := FetchChangelog(context.Background(), client, DefaultChangelogURL)
	if err != nil {
		t.Fatalf("FetchChangelog returned error: %v", err)
	}

	if string(got) != body {
		t.Fatalf("unexpected body: %q", string(got))
	}
}

func TestFetchChangelog_HTTPError(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("boom")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	_, err := FetchChangelog(context.Background(), client, DefaultChangelogURL)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status") {
		t.Fatalf("expected status error, got %v", err)
	}
}

func TestFetchChangelog_RejectsDisallowedHost(t *testing.T) {
	t.Parallel()

	_, err := FetchChangelog(context.Background(), nil, "https://example.invalid/changelog")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "host") {
		t.Fatalf("expected host validation error, got %v", err)
	}
}

func TestFilterFeatureReleases(t *testing.T) {
	t.Parallel()

	releases := []Release{
		{
			Version: "v1.54.0",
			Notes: []string{
				"**Feature**: Adding new attribute to disable IMDS v1 APIs for fleet, Image Builder and AppBlockBuilder instances.",
			},
		},
		{
			Version: "v1.53.2",
			Notes: []string{
				"No change notes available for this release.",
			},
		},
		{
			Version: "v1.53.1",
			Notes: []string{
				"**Dependency Update**: Updated to the latest SDK module versions",
			},
		},
		{
			Version: "v1.53.0",
			Notes: []string{
				"**Dependency Update**: Updated to the latest SDK module versions",
				"**Feature**: Added support for new operating systems",
			},
		},
		{
			Version: "v1.52.0",
			Notes: []string{
				"**Feature**: Adding support for additional instances and extended storage",
			},
		},
	}

	got := FilterFeatureReleases(releases)

	if len(got) != 3 {
		t.Fatalf("expected 3 feature releases, got %d", len(got))
	}

	if got[0].Version != "v1.54.0" {
		t.Fatalf("expected first feature release v1.54.0, got %q", got[0].Version)
	}

	if got[1].Version != "v1.53.0" {
		t.Fatalf("expected second feature release v1.53.0, got %q", got[1].Version)
	}

	if got[2].Version != "v1.52.0" {
		t.Fatalf("expected third feature release v1.52.0, got %q", got[2].Version)
	}
}

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
