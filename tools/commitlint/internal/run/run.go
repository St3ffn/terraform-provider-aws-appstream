// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package run

import (
	"context"
	"fmt"

	"github.com/st3ffn/terraform-provider-aws-appstream/tools/commitlint/internal/conventional"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/commitlint/internal/gitrepo"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/commitlint/internal/model"
)

// Lint selects commits according to options and validates each selected commit
// against Conventional Commits.
func Lint(ctx context.Context, opts Options) (model.Report, error) {
	if err := ctx.Err(); err != nil {
		return model.Report{}, err
	}

	opts = applyDefaults(opts)
	if err := validateOptions(opts); err != nil {
		return model.Report{}, err
	}

	repo, err := gitrepo.Open(opts.RepoPath)
	if err != nil {
		return model.Report{}, err
	}

	commits, err := selectCommits(ctx, repo, opts)
	if err != nil {
		return model.Report{}, err
	}

	report := model.Report{TotalCommits: len(commits)}
	for _, commit := range commits {
		if err := ctx.Err(); err != nil {
			return model.Report{}, err
		}

		if !opts.IncludeMergeCommits && conventional.IsMergeCommitSubject(commit.Subject) {
			report.SkippedCommits++
			continue
		}

		violations := conventional.ValidateCommitMessage(commit.Message)
		result := model.CommitResult{
			Hash:      commit.Hash,
			ShortHash: commit.ShortHash,
			Subject:   commit.Subject,
			Message:   commit.Message,
			Valid:     len(violations) == 0,
		}
		if len(violations) > 0 {
			result.Violations = violations
			report.InvalidCommits++
		}

		report.CheckedCommits++
		report.Results = append(report.Results, result)
	}

	return report, nil
}

func selectCommits(ctx context.Context, repo *gitrepo.Repository, opts Options) ([]gitrepo.Commit, error) {
	if opts.All {
		return repo.SelectAll(ctx)
	}

	if opts.Branch != "" {
		return repo.SelectBranch(ctx, opts.Branch)
	}

	if opts.BranchDiff != "" {
		commits, err := repo.SelectBranchDiff(ctx, opts.BranchDiff, opts.Base)
		if err == nil {
			return commits, nil
		}
		if opts.Base == defaultBaseRef {
			fallback, fallbackErr := repo.SelectBranchDiff(ctx, opts.BranchDiff, "origin/HEAD")
			if fallbackErr == nil {
				return fallback, nil
			}
		}
		return nil, err
	}

	if opts.From != "" && opts.To != "" {
		return repo.SelectRange(ctx, opts.From, opts.To)
	}

	commits, err := repo.SelectCurrentBranchDiff(ctx, opts.Base)
	if err == nil {
		return commits, nil
	}
	if opts.Base == defaultBaseRef {
		fallback, fallbackErr := repo.SelectCurrentBranchDiff(ctx, "origin/HEAD")
		if fallbackErr == nil {
			return fallback, nil
		}
	}
	return nil, fmt.Errorf("select default commit set: %w", err)
}
