// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package run

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/githubissues"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/appstream-changelog-issues/internal/report"
)

func Execute(ctx context.Context, opts Options) error {
	return execute(ctx, opts, defaultDependencies())
}

func execute(ctx context.Context, opts Options, deps dependencies) error {
	if err := validateOptions(opts); err != nil {
		return err
	}
	opts = applyDefaults(opts)
	log.Printf("start repo=%s/%s dry_run=%t go_mod=%s", opts.RepoOwner, opts.RepoName, opts.DryRun, opts.ProviderGoModPath)

	currentVersion, err := deps.appStreamVersion(opts.ProviderGoModPath)
	if err != nil {
		return fmt.Errorf("read current appstream sdk version: %w", err)
	}
	log.Printf("current_version=%s", currentVersion)

	changelog, err := deps.fetchChangelog(ctx, nil, opts.ChangelogURL)
	if err != nil {
		return fmt.Errorf("fetch appstream changelog: %w", err)
	}

	releases, err := deps.parseChangelog(changelog)
	if err != nil {
		return fmt.Errorf("parse appstream changelog: %w", err)
	}
	log.Printf("parsed_releases=%d", len(releases))

	newerReleases, err := deps.newerThan(releases, currentVersion)
	if err != nil {
		return fmt.Errorf("filter releases newer than %s: %w", currentVersion, err)
	}
	log.Printf("newer_releases=%d", len(newerReleases))
	if len(newerReleases) == 0 {
		log.Printf("no newer releases, exiting")
		return nil
	}

	featureReleases := deps.filterFeatureRelease(newerReleases)
	log.Printf("feature_releases=%d", len(featureReleases))
	if len(featureReleases) == 0 {
		log.Printf("no feature releases, exiting")
		return nil
	}

	issueContent, err := deps.buildIssueContent(currentVersion, featureReleases, report.BuildOptions{
		IssueTitle:  opts.IssueTitle,
		IssueMarker: opts.IssueMarker,
		SourceURL:   opts.ChangelogURL,
	})
	if err != nil {
		return fmt.Errorf("build issue content: %w", err)
	}

	if opts.DryRun {
		log.Printf("dry-run: would upsert issue title=%q labels=%d", issueContent.Title, len(opts.IssueLabels))
		return nil
	}
	if opts.GitHubToken == nil || !opts.GitHubToken.IsSet() {
		return errors.New("github token must be set when dry-run is false")
	}

	issuesClient, err := deps.newGitHubClient(opts.RepoOwner, opts.RepoName, opts.GitHubToken)
	if err != nil {
		return fmt.Errorf("create github issues client: %w", err)
	}

	result, err := issuesClient.CreateOrUpdateIssue(ctx, githubissues.UpsertInput{
		Title:  issueContent.Title,
		Body:   issueContent.Body,
		Marker: opts.IssueMarker,
		Labels: opts.IssueLabels,
	})
	if err != nil {
		return fmt.Errorf("create or update github issue: %w", err)
	}
	log.Printf("issue_%s number=%d url=%s", result.Action, result.IssueNumber, result.IssueURL)

	return nil
}
