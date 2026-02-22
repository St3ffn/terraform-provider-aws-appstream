// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/githubissues"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/run"
)

var (
	providerGoMod = flag.String("provider-go-mod", "go.mod", "Path to provider go.mod")
	repo          = flag.String("repo", "", "GitHub repository in owner/repo format (defaults to GITHUB_REPOSITORY)")
	tokenFlag     = flag.String("github-token", "", "GitHub token (defaults to GITHUB_TOKEN)")
	dryRun        = flag.Bool("dry-run", false, "Compute changes without creating/updating issues")
	timeout       = flag.Duration("timeout", 5*time.Minute, "Overall execution timeout")

	stderr = os.Stderr
)

func main() {
	log.SetFlags(0)
	os.Exit(runCLI())
}

func runCLI() int {
	flag.Parse()

	repoVal := strings.TrimSpace(*repo)
	if repoVal == "" {
		repoVal = strings.TrimSpace(os.Getenv("GITHUB_REPOSITORY"))
	}
	owner, name, ok := strings.Cut(repoVal, "/")
	if !ok || owner == "" || name == "" {
		_, _ = fmt.Fprintf(stderr, "invalid repo, expected owner/repo, got %q\n", repoVal)
		return 2
	}

	token := strings.TrimSpace(*tokenFlag)
	if token == "" {
		token = strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	}
	if token == "" && !*dryRun {
		_, _ = fmt.Fprintln(stderr, "missing GitHub token (set GITHUB_TOKEN or --github-token)")
		return 2
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	err := run.Execute(ctx, run.Options{
		ProviderGoModPath: *providerGoMod,
		RepoOwner:         owner,
		RepoName:          name,
		GitHubToken:       githubissues.NewToken(token),
		DryRun:            *dryRun,
	})
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "appstream-changelog-issues: %v\n", err)
		return 1
	}

	return 0
}
