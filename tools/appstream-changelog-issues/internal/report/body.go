// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package report

import (
	"errors"
	"fmt"
	"strings"

	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/appstream"
	"golang.org/x/mod/semver"
)

const (
	// IssueTitle is the stable issue title used for changelog watch upserts.
	IssueTitle = "AWS SDK service/appstream feature updates available"
	// IssueMarker is a hidden marker used to detect/update an existing tracking issue.
	IssueMarker = "<!-- appstream-changelog-issues -->"
	// DefaultSourceURL points to the source changelog.
	DefaultSourceURL = appstream.DefaultChangelogURL
)

// IssueContent contains title and body payload for GitHub issue create/update.
type IssueContent struct {
	Title string
	Body  string
}

// BuildOptions controls rendering of the generated tracking issue content.
type BuildOptions struct {
	IssueTitle  string
	IssueMarker string
	SourceURL   string
}

// BuildIssueContent renders a deterministic issue payload from feature releases.
func BuildIssueContent(currentVersion string, releases []appstream.Release, sourceURL string) (IssueContent, error) {
	return BuildIssueContentWithOptions(currentVersion, releases, BuildOptions{
		IssueTitle:  IssueTitle,
		IssueMarker: IssueMarker,
		SourceURL:   sourceURL,
	})
}

// BuildIssueContentWithOptions renders a deterministic issue payload from feature
// releases using caller-provided title/marker/source values.
func BuildIssueContentWithOptions(currentVersion string, releases []appstream.Release, opts BuildOptions) (IssueContent, error) {
	currentVersion = strings.TrimSpace(currentVersion)
	if !semver.IsValid(currentVersion) {
		return IssueContent{}, fmt.Errorf("invalid current version %q", currentVersion)
	}

	if len(releases) == 0 {
		return IssueContent{}, errors.New("releases must not be empty")
	}

	issueTitle := strings.TrimSpace(opts.IssueTitle)
	if issueTitle == "" {
		issueTitle = IssueTitle
	}
	issueMarker := strings.TrimSpace(opts.IssueMarker)
	if issueMarker == "" {
		issueMarker = IssueMarker
	}
	sourceURL := strings.TrimSpace(opts.SourceURL)
	if sourceURL == "" {
		sourceURL = DefaultSourceURL
	}

	latest, err := latestVersion(releases)
	if err != nil {
		return IssueContent{}, err
	}

	var b strings.Builder
	b.WriteString(issueMarker)
	b.WriteString("\n")
	b.WriteString("## Summary\n")
	b.WriteString(fmt.Sprintf("- Current `github.com/aws/aws-sdk-go-v2/service/appstream` version: `%s`\n", currentVersion))
	b.WriteString(fmt.Sprintf("- Latest feature release detected: `%s`\n", latest))
	b.WriteString(fmt.Sprintf("- Feature releases found: `%d`\n", len(releases)))
	b.WriteString("\n")
	b.WriteString("## Feature Releases\n")

	for _, release := range releases {
		b.WriteString(fmt.Sprintf("### %s (%s)\n", release.Version, release.Date))
		for _, note := range release.Notes {
			if strings.TrimSpace(note) == "" {
				continue
			}
			b.WriteString("- ")
			b.WriteString(note)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("## Source\n")
	b.WriteString(fmt.Sprintf("- %s\n", sourceURL))

	return IssueContent{
		Title: issueTitle,
		Body:  b.String(),
	}, nil
}

func latestVersion(releases []appstream.Release) (string, error) {
	var latest string

	for _, release := range releases {
		if !semver.IsValid(release.Version) {
			return "", fmt.Errorf("invalid release version %q", release.Version)
		}

		if latest == "" || semver.Compare(release.Version, latest) > 0 {
			latest = release.Version
		}
	}

	if latest == "" {
		return "", errors.New("no valid release versions found")
	}

	return latest, nil
}
