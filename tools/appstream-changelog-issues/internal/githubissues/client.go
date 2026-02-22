// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package githubissues

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/v83/github"
)

// Client wraps GitHub issue operations for a single repository.
type Client struct {
	owner  string
	repo   string
	issues issuesService
}

type issuesService interface {
	ListByRepo(ctx context.Context, owner, repo string, opts *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error)
	Create(ctx context.Context, owner, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
	Edit(ctx context.Context, owner, repo string, number int, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
}

// NewClient creates a GitHub issue client for a fixed owner/repository pair.
func NewClient(owner, repo string, token *Token) (*Client, error) {
	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)

	if owner == "" {
		return nil, errors.New("owner must not be empty")
	}
	if repo == "" {
		return nil, errors.New("repo must not be empty")
	}
	if token == nil || !token.IsSet() {
		return nil, errors.New("token must not be empty")
	}

	gh := github.NewClient(nil).WithAuthToken(token.value())

	return &Client{
		owner:  owner,
		repo:   repo,
		issues: gh.Issues,
	}, nil
}

func (c *Client) validate() error {
	if c == nil {
		return errors.New("client must not be nil")
	}
	if strings.TrimSpace(c.owner) == "" || strings.TrimSpace(c.repo) == "" {
		return errors.New("client owner/repo must not be empty")
	}
	if c.issues == nil {
		return errors.New("client issues service must not be nil")
	}
	return nil
}

func wrapGitHubError(op string, err error) error {
	return fmt.Errorf("github issues %s: %w", op, err)
}
