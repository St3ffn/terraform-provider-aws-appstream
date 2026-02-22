// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package run

import (
	"context"
	"net/http"

	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/appstream"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/githubissues"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/gomod"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/report"
)

type githubIssueClient interface {
	CreateOrUpdateIssue(ctx context.Context, input githubissues.UpsertInput) (githubissues.UpsertResult, error)
}

type dependencies struct {
	appStreamVersion     func(goModPath string) (string, error)
	fetchChangelog       func(ctx context.Context, client *http.Client, url string) ([]byte, error)
	parseChangelog       func(markdown []byte) ([]appstream.Release, error)
	newerThan            func(releases []appstream.Release, currentVersion string) ([]appstream.Release, error)
	filterFeatureRelease func(releases []appstream.Release) []appstream.Release
	buildIssueContent    func(currentVersion string, releases []appstream.Release, opts report.BuildOptions) (report.IssueContent, error)
	newGitHubClient      func(owner, repo string, token *githubissues.Token) (githubIssueClient, error)
}

func defaultDependencies() dependencies {
	return dependencies{
		appStreamVersion:     gomod.AppStreamVersion,
		fetchChangelog:       appstream.FetchChangelog,
		parseChangelog:       appstream.ParseChangelog,
		newerThan:            appstream.NewerThan,
		filterFeatureRelease: appstream.FilterFeatureReleases,
		buildIssueContent:    report.BuildIssueContentWithOptions,
		newGitHubClient: func(owner, repo string, token *githubissues.Token) (githubIssueClient, error) {
			return githubissues.NewClient(owner, repo, token)
		},
	}
}
