// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package githubissues

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-github/v83/github"
)

type UpsertAction string

const (
	UpsertActionCreated UpsertAction = "created"
	UpsertActionUpdated UpsertAction = "updated"
	UpsertActionNoop    UpsertAction = "noop"
)

type UpsertInput struct {
	Title  string
	Body   string
	Marker string
	Labels []string
}

type UpsertResult struct {
	Action      UpsertAction
	IssueNumber int
	IssueURL    string
}

// CreateOrUpdateIssue creates or updates an issue for the configured repository.
//
// Matching logic:
// 1) open issue body contains Marker
// 2) fallback: open issue title matches Title
func (c *Client) CreateOrUpdateIssue(ctx context.Context, input UpsertInput) (UpsertResult, error) {
	if err := c.validate(); err != nil {
		return UpsertResult{}, err
	}
	if err := validateUpsertInput(input); err != nil {
		return UpsertResult{}, err
	}

	desiredLabels := normalizeLabels(input.Labels)
	targetIssue, err := c.findTargetIssue(ctx, input.Title, input.Marker)
	if err != nil {
		return UpsertResult{}, err
	}

	if targetIssue == nil {
		req := &github.IssueRequest{
			Title: github.Ptr(input.Title),
			Body:  github.Ptr(input.Body),
		}
		if len(desiredLabels) > 0 {
			req.Labels = &desiredLabels
		}

		created, _, err := c.issues.Create(ctx, c.owner, c.repo, req)
		if err != nil {
			return UpsertResult{}, wrapGitHubError("create", err)
		}
		if created == nil {
			return UpsertResult{}, errors.New("github issues create returned nil issue")
		}

		return UpsertResult{
			Action:      UpsertActionCreated,
			IssueNumber: created.GetNumber(),
			IssueURL:    created.GetHTMLURL(),
		}, nil
	}

	needsTitleUpdate := targetIssue.GetTitle() != input.Title
	needsBodyUpdate := targetIssue.GetBody() != input.Body
	missing := missingLabels(targetIssue, desiredLabels)

	if !needsTitleUpdate && !needsBodyUpdate && len(missing) == 0 {
		return UpsertResult{
			Action:      UpsertActionNoop,
			IssueNumber: targetIssue.GetNumber(),
			IssueURL:    targetIssue.GetHTMLURL(),
		}, nil
	}

	req := &github.IssueRequest{}
	if needsTitleUpdate {
		req.Title = github.Ptr(input.Title)
	}
	if needsBodyUpdate {
		req.Body = github.Ptr(input.Body)
	}
	if len(missing) > 0 {
		merged := unionLabels(issueLabelNames(targetIssue), desiredLabels)
		req.Labels = &merged
	}

	updated, _, err := c.issues.Edit(ctx, c.owner, c.repo, targetIssue.GetNumber(), req)
	if err != nil {
		return UpsertResult{}, wrapGitHubError("edit", err)
	}
	if updated == nil {
		return UpsertResult{}, errors.New("github issues edit returned nil issue")
	}

	return UpsertResult{
		Action:      UpsertActionUpdated,
		IssueNumber: updated.GetNumber(),
		IssueURL:    updated.GetHTMLURL(),
	}, nil
}

func validateUpsertInput(input UpsertInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return errors.New("upsert input title must not be empty")
	}
	if strings.TrimSpace(input.Body) == "" {
		return errors.New("upsert input body must not be empty")
	}
	if strings.TrimSpace(input.Marker) == "" {
		return errors.New("upsert input marker must not be empty")
	}
	return nil
}

func (c *Client) findTargetIssue(ctx context.Context, title, marker string) (*github.Issue, error) {
	markerMatches := make([]*github.Issue, 0, 1)
	titleMatches := make([]*github.Issue, 0, 1)

	opts := &github.IssueListByRepoOptions{
		State:       "open",
		Sort:        "created",
		Direction:   "asc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		issues, resp, err := c.issues.ListByRepo(ctx, c.owner, c.repo, opts)
		if err != nil {
			return nil, wrapGitHubError("list", err)
		}

		for _, issue := range issues {
			if issue == nil || issue.PullRequestLinks != nil {
				continue
			}

			if strings.Contains(issue.GetBody(), marker) {
				markerMatches = append(markerMatches, issue)
				continue
			}
			if issue.GetTitle() == title {
				titleMatches = append(titleMatches, issue)
			}
		}

		if resp == nil || resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}

	switch len(markerMatches) {
	case 0:
		// continue with title fallback
	case 1:
		return markerMatches[0], nil
	default:
		return nil, fmt.Errorf("found multiple open issues with marker %q", marker)
	}

	switch len(titleMatches) {
	case 0:
		return nil, nil
	case 1:
		return titleMatches[0], nil
	default:
		return nil, fmt.Errorf("found multiple open issues with title %q", title)
	}
}

func issueLabelNames(issue *github.Issue) []string {
	if issue == nil {
		return nil
	}

	out := make([]string, 0, len(issue.Labels))
	for _, label := range issue.Labels {
		if label == nil {
			continue
		}
		name := strings.TrimSpace(label.GetName())
		if name == "" {
			continue
		}
		out = append(out, name)
	}
	return normalizeLabels(out)
}

func missingLabels(issue *github.Issue, desired []string) []string {
	if len(desired) == 0 {
		return nil
	}

	current := issueLabelNames(issue)
	missing := make([]string, 0, len(desired))
	for _, label := range desired {
		if !slices.Contains(current, label) {
			missing = append(missing, label)
		}
	}
	return missing
}

func unionLabels(current, desired []string) []string {
	merged := make([]string, 0, len(current)+len(desired))
	merged = append(merged, current...)
	merged = append(merged, desired...)
	return normalizeLabels(merged)
}

func normalizeLabels(labels []string) []string {
	if len(labels) == 0 {
		return nil
	}

	out := make([]string, 0, len(labels))
	for _, label := range labels {
		label = strings.TrimSpace(label)
		if label == "" {
			continue
		}
		out = append(out, label)
	}

	if len(out) == 0 {
		return nil
	}

	slices.Sort(out)
	return slices.Compact(out)
}
