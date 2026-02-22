// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package run

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/appstream"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/githubissues"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/report"
)

func TestExecute_DryRunSkipsGitHubIssueWrite(t *testing.T) {
	t.Parallel()

	clientCreated := false
	deps := fixedDeps()
	deps.newGitHubClient = func(owner, repo string, token *githubissues.Token) (githubIssueClient, error) {
		clientCreated = true
		return &fakeGitHubIssueClient{}, nil
	}

	err := execute(context.Background(), Options{
		ProviderGoModPath: "go.mod",
		RepoOwner:         "owner",
		RepoName:          "repo",
		GitHubToken:       githubissues.NewToken("token"),
		DryRun:            true,
	}, deps)
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if clientCreated {
		t.Fatal("did not expect github client to be created in dry-run mode")
	}
}

func TestExecute_NoNewerReleasesIsNoop(t *testing.T) {
	t.Parallel()

	clientCreated := false
	deps := fixedDeps()
	deps.newerThan = func(_ []appstream.Release, _ string) ([]appstream.Release, error) {
		return nil, nil
	}
	deps.newGitHubClient = func(owner, repo string, token *githubissues.Token) (githubIssueClient, error) {
		clientCreated = true
		return &fakeGitHubIssueClient{}, nil
	}

	err := execute(context.Background(), Options{
		ProviderGoModPath: "go.mod",
		RepoOwner:         "owner",
		RepoName:          "repo",
		GitHubToken:       githubissues.NewToken("token"),
	}, deps)
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if clientCreated {
		t.Fatal("did not expect github client when there are no newer releases")
	}
}

func TestExecute_NonFeatureReleasesIsNoop(t *testing.T) {
	t.Parallel()

	clientCreated := false
	deps := fixedDeps()
	deps.filterFeatureRelease = func(_ []appstream.Release) []appstream.Release {
		return nil
	}
	deps.newGitHubClient = func(owner, repo string, token *githubissues.Token) (githubIssueClient, error) {
		clientCreated = true
		return &fakeGitHubIssueClient{}, nil
	}

	err := execute(context.Background(), Options{
		ProviderGoModPath: "go.mod",
		RepoOwner:         "owner",
		RepoName:          "repo",
		GitHubToken:       githubissues.NewToken("token"),
	}, deps)
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if clientCreated {
		t.Fatal("did not expect github client when there are no feature releases")
	}
}

func TestExecute_RequiresTokenWhenNotDryRun(t *testing.T) {
	t.Parallel()

	err := execute(context.Background(), Options{
		ProviderGoModPath: "go.mod",
		RepoOwner:         "owner",
		RepoName:          "repo",
		DryRun:            false,
	}, fixedDeps())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "github token must be set") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_CreateOrUpdateIssue(t *testing.T) {
	t.Parallel()

	fakeClient := &fakeGitHubIssueClient{
		result: githubissues.UpsertResult{
			Action:      githubissues.UpsertActionUpdated,
			IssueNumber: 123,
			IssueURL:    "https://example.invalid/issues/123",
		},
	}

	deps := fixedDeps()
	deps.newGitHubClient = func(owner, repo string, token *githubissues.Token) (githubIssueClient, error) {
		if owner != "owner" || repo != "repo" {
			t.Fatalf("unexpected owner/repo: %s/%s", owner, repo)
		}
		if token == nil || !token.IsSet() {
			t.Fatal("expected non-empty token")
		}
		return fakeClient, nil
	}

	err := execute(context.Background(), Options{
		ProviderGoModPath: "go.mod",
		RepoOwner:         "owner",
		RepoName:          "repo",
		GitHubToken:       githubissues.NewToken("token"),
	}, deps)
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if !fakeClient.called {
		t.Fatal("expected CreateOrUpdateIssue to be called")
	}
	if fakeClient.input.Title != report.IssueTitle {
		t.Fatalf("expected issue title %q, got %q", report.IssueTitle, fakeClient.input.Title)
	}
	if fakeClient.input.Marker != report.IssueMarker {
		t.Fatalf("expected marker %q, got %q", report.IssueMarker, fakeClient.input.Marker)
	}
	if !slices.Equal(fakeClient.input.Labels, defaultIssueLabels) {
		t.Fatalf("expected labels %v, got %v", defaultIssueLabels, fakeClient.input.Labels)
	}
}

func TestExecute_UsesCustomIssueOptions(t *testing.T) {
	t.Parallel()

	fakeClient := &fakeGitHubIssueClient{}
	deps := fixedDeps()
	deps.buildIssueContent = func(currentVersion string, releases []appstream.Release, opts report.BuildOptions) (report.IssueContent, error) {
		if opts.IssueTitle != "Custom title" {
			t.Fatalf("expected custom title, got %q", opts.IssueTitle)
		}
		if opts.IssueMarker != "<!-- custom -->" {
			t.Fatalf("expected custom marker, got %q", opts.IssueMarker)
		}
		if opts.SourceURL != "https://example.invalid/changelog" {
			t.Fatalf("expected custom changelog url, got %q", opts.SourceURL)
		}
		return report.IssueContent{
			Title: "Custom title",
			Body:  "<!-- custom -->\ncontent",
		}, nil
	}
	deps.newGitHubClient = func(owner, repo string, token *githubissues.Token) (githubIssueClient, error) {
		return fakeClient, nil
	}

	err := execute(context.Background(), Options{
		ProviderGoModPath: "go.mod",
		RepoOwner:         "owner",
		RepoName:          "repo",
		GitHubToken:       githubissues.NewToken("token"),
		ChangelogURL:      "https://example.invalid/changelog",
		IssueTitle:        "Custom title",
		IssueMarker:       "<!-- custom -->",
		IssueLabels:       []string{"x", "y"},
	}, deps)
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if fakeClient.input.Title != "Custom title" {
		t.Fatalf("expected custom issue title, got %q", fakeClient.input.Title)
	}
	if fakeClient.input.Marker != "<!-- custom -->" {
		t.Fatalf("expected custom marker, got %q", fakeClient.input.Marker)
	}
	if !slices.Equal(fakeClient.input.Labels, []string{"x", "y"}) {
		t.Fatalf("expected custom labels, got %v", fakeClient.input.Labels)
	}
}

func TestExecute_ValidateOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts Options
	}{
		{
			name: "missing go.mod path",
			opts: Options{
				RepoOwner: "owner",
				RepoName:  "repo",
			},
		},
		{
			name: "missing owner",
			opts: Options{
				ProviderGoModPath: "go.mod",
				RepoName:          "repo",
			},
		},
		{
			name: "missing repo",
			opts: Options{
				ProviderGoModPath: "go.mod",
				RepoOwner:         "owner",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := execute(context.Background(), tt.opts, fixedDeps())
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestExecute_PropagatesCreateOrUpdateError(t *testing.T) {
	t.Parallel()

	deps := fixedDeps()
	deps.newGitHubClient = func(owner, repo string, token *githubissues.Token) (githubIssueClient, error) {
		return &fakeGitHubIssueClient{err: errors.New("boom")}, nil
	}

	err := execute(context.Background(), Options{
		ProviderGoModPath: "go.mod",
		RepoOwner:         "owner",
		RepoName:          "repo",
		GitHubToken:       githubissues.NewToken("token"),
	}, deps)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "create or update github issue") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func fixedDeps() dependencies {
	releases := []appstream.Release{
		{
			Version: "v1.54.0",
			Date:    "2026-02-20",
			Notes:   []string{"**Feature**: Added support for X"},
		},
	}

	return dependencies{
		appStreamVersion: func(goModPath string) (string, error) {
			if strings.TrimSpace(goModPath) == "" {
				return "", errors.New("go.mod path is empty")
			}
			return "v1.53.0", nil
		},
		fetchChangelog: func(ctx context.Context, client *http.Client, url string) ([]byte, error) {
			return []byte("dummy"), nil
		},
		parseChangelog: func(markdown []byte) ([]appstream.Release, error) {
			return releases, nil
		},
		newerThan: func(releases []appstream.Release, currentVersion string) ([]appstream.Release, error) {
			return releases, nil
		},
		filterFeatureRelease: func(releases []appstream.Release) []appstream.Release {
			return releases
		},
		buildIssueContent: func(currentVersion string, releases []appstream.Release, opts report.BuildOptions) (report.IssueContent, error) {
			return report.IssueContent{
				Title: opts.IssueTitle,
				Body:  opts.IssueMarker + "\ncontent",
			}, nil
		},
		newGitHubClient: func(owner, repo string, token *githubissues.Token) (githubIssueClient, error) {
			return &fakeGitHubIssueClient{}, nil
		},
	}
}

type fakeGitHubIssueClient struct {
	called bool
	input  githubissues.UpsertInput
	result githubissues.UpsertResult
	err    error
}

func (f *fakeGitHubIssueClient) CreateOrUpdateIssue(_ context.Context, input githubissues.UpsertInput) (githubissues.UpsertResult, error) {
	f.called = true
	f.input = input
	if f.err != nil {
		return githubissues.UpsertResult{}, f.err
	}
	return f.result, nil
}
