// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package run

import (
	"errors"
	"strings"

	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/appstream"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/githubissues"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/report"
)

const (
	defaultIssueLabelDependencies = "dependencies"
	defaultIssueLabelWatch        = "appstream-sdk-watch"
)

var defaultIssueLabels = []string{
	defaultIssueLabelDependencies,
	defaultIssueLabelWatch,
}

type Options struct {
	ProviderGoModPath string
	RepoOwner         string
	RepoName          string
	GitHubToken       *githubissues.Token
	ChangelogURL      string
	IssueTitle        string
	IssueMarker       string
	IssueLabels       []string
	DryRun            bool
}

func applyDefaults(opts Options) Options {
	if strings.TrimSpace(opts.ChangelogURL) == "" {
		opts.ChangelogURL = appstream.DefaultChangelogURL
	}
	if strings.TrimSpace(opts.IssueTitle) == "" {
		opts.IssueTitle = report.IssueTitle
	}
	if strings.TrimSpace(opts.IssueMarker) == "" {
		opts.IssueMarker = report.IssueMarker
	}
	if len(opts.IssueLabels) == 0 {
		opts.IssueLabels = append([]string(nil), defaultIssueLabels...)
	}
	return opts
}

func validateOptions(opts Options) error {
	if strings.TrimSpace(opts.ProviderGoModPath) == "" {
		return errors.New("provider go.mod path must not be empty")
	}
	if strings.TrimSpace(opts.RepoOwner) == "" {
		return errors.New("repo owner must not be empty")
	}
	if strings.TrimSpace(opts.RepoName) == "" {
		return errors.New("repo name must not be empty")
	}
	return nil
}
