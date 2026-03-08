// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package gitrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// Commit represents a commit selected for linting.
type Commit struct {
	Hash      string
	ShortHash string
	Subject   string
	Message   string
}

// Repository wraps a go-git repository and exposes selection helpers.
type Repository struct {
	repo *git.Repository
}

// Open opens a git repository from path, detecting .git directories.
func Open(path string) (*Repository, error) {
	r, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, fmt.Errorf("open repository %q: %w", path, err)
	}
	return &Repository{repo: r}, nil
}

// SelectAll returns all commits reachable from all refs.
func (r *Repository) SelectAll(ctx context.Context) ([]Commit, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	iter, err := r.repo.Log(&git.LogOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("list all commits: %w", err)
	}
	defer iter.Close()

	seen := map[plumbing.Hash]struct{}{}
	var out []Commit
	err = iter.ForEach(func(c *object.Commit) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, ok := seen[c.Hash]; ok {
			return nil
		}
		seen[c.Hash] = struct{}{}
		out = append(out, toCommit(c))
		return nil
	})
	if err != nil {
		if isObjectNotFoundError(err) && len(out) > 0 {
			return out, nil
		}
		return nil, fmt.Errorf("iterate all commits: %w", err)
	}

	return out, nil
}

// SelectBranch returns all commits reachable from the branch tip.
func (r *Repository) SelectBranch(ctx context.Context, branch string) ([]Commit, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	tip, err := r.resolveBranchTip(branch)
	if err != nil {
		return nil, err
	}

	return r.commitsFrom(ctx, tip)
}

// SelectBranchDiff returns commits reachable from branch but not reachable from base.
func (r *Repository) SelectBranchDiff(ctx context.Context, branch, base string) ([]Commit, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	branchTip, err := r.resolveBranchTip(branch)
	if err != nil {
		return nil, fmt.Errorf("resolve branch %q: %w", branch, err)
	}
	baseTip, err := r.resolveBranchTip(base)
	if err != nil {
		return nil, fmt.Errorf("resolve base %q: %w", base, err)
	}

	return r.diffFromTips(ctx, branchTip, baseTip)
}

// SelectCurrentBranchDiff returns commits reachable from HEAD but not from base.
func (r *Repository) SelectCurrentBranchDiff(ctx context.Context, base string) ([]Commit, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	head, err := r.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("resolve HEAD: %w", err)
	}

	baseTip, err := r.resolveBranchTip(base)
	if err != nil {
		return nil, fmt.Errorf("resolve base %q: %w", base, err)
	}

	return r.diffFromTips(ctx, head.Hash(), baseTip)
}

// SelectRange returns commits in an inclusive range from..to.
//
// It follows git-style semantics similar to from^..to by including commits
// reachable from to and excluding ancestors of all parents of from.
func (r *Repository) SelectRange(ctx context.Context, from, to string) ([]Commit, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	fromHash, err := r.resolveRevision(from)
	if err != nil {
		return nil, fmt.Errorf("resolve from %q: %w", from, err)
	}
	toHash, err := r.resolveRevision(to)
	if err != nil {
		return nil, fmt.Errorf("resolve to %q: %w", to, err)
	}

	toReach, err := r.reachableSet(ctx, toHash)
	if err != nil {
		if isObjectNotFoundError(err) {
			return r.selectRangeLinear(ctx, fromHash, toHash)
		}
		return nil, fmt.Errorf("compute to reachability: %w", err)
	}
	if _, ok := toReach[fromHash]; !ok {
		return nil, fmt.Errorf("from commit %s is not reachable from to commit %s", fromHash.String(), toHash.String())
	}

	fromCommit, err := r.repo.CommitObject(fromHash)
	if err != nil {
		if isObjectNotFoundError(err) {
			return r.selectRangeLinear(ctx, fromHash, toHash)
		}
		return nil, fmt.Errorf("load from commit %s: %w", fromHash.String(), err)
	}

	exclude := map[plumbing.Hash]struct{}{}
	for _, parent := range fromCommit.ParentHashes {
		parentReach, err := r.reachableSet(ctx, parent)
		if err != nil {
			if isObjectNotFoundError(err) {
				return r.selectRangeLinear(ctx, fromHash, toHash)
			}
			return nil, fmt.Errorf("compute parent reachability: %w", err)
		}
		for h := range parentReach {
			exclude[h] = struct{}{}
		}
	}

	iter, err := r.repo.Log(&git.LogOptions{From: toHash})
	if err != nil {
		return nil, fmt.Errorf("list range commits: %w", err)
	}
	defer iter.Close()

	var out []Commit
	seen := map[plumbing.Hash]struct{}{}
	err = iter.ForEach(func(c *object.Commit) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, drop := exclude[c.Hash]; drop {
			return nil
		}
		if _, ok := seen[c.Hash]; ok {
			return nil
		}
		seen[c.Hash] = struct{}{}
		out = append(out, toCommit(c))
		return nil
	})
	if err != nil {
		if isObjectNotFoundError(err) {
			return r.selectRangeLinear(ctx, fromHash, toHash)
		}
		return nil, fmt.Errorf("iterate range commits: %w", err)
	}

	return out, nil
}

func (r *Repository) selectRangeLinear(ctx context.Context, fromHash, toHash plumbing.Hash) ([]Commit, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	iter, err := r.repo.Log(&git.LogOptions{From: toHash})
	if err != nil {
		return nil, fmt.Errorf("list range commits: %w", err)
	}
	defer iter.Close()

	var (
		out   []Commit
		found bool
	)
	err = iter.ForEach(func(c *object.Commit) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		out = append(out, toCommit(c))
		if c.Hash == fromHash {
			found = true
			return storer.ErrStop
		}
		return nil
	})
	if err != nil && !errors.Is(err, storer.ErrStop) {
		if isObjectNotFoundError(err) && found {
			return out, nil
		}
		return nil, fmt.Errorf("iterate range commits: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("from commit %s is not reachable from to commit %s", fromHash.String(), toHash.String())
	}

	return out, nil
}

func (r *Repository) diffFromTips(ctx context.Context, branchTip, baseTip plumbing.Hash) ([]Commit, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// If both tips are identical, the diff set is empty and no graph walk is needed.
	// This avoids unnecessary traversal failures in repos with limited history.
	if branchTip == baseTip {
		return []Commit{}, nil
	}

	baseReach, err := r.reachableSet(ctx, baseTip)
	if err != nil {
		// Some repositories can contain refs that are valid for git CLI traversal
		// but fail in go-git history walks with "object not found". In this case
		// fall back to a linear walk from branch tip to the base commit.
		if isObjectNotFoundError(err) {
			return r.diffByLinearHistory(ctx, branchTip, baseTip)
		}
		return nil, fmt.Errorf("compute base reachability: %w", err)
	}

	iter, err := r.repo.Log(&git.LogOptions{From: branchTip})
	if err != nil {
		return nil, fmt.Errorf("list branch commits: %w", err)
	}
	defer iter.Close()

	var out []Commit
	err = iter.ForEach(func(c *object.Commit) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, existsInBase := baseReach[c.Hash]; existsInBase {
			return nil
		}
		out = append(out, toCommit(c))
		return nil
	})
	if err != nil {
		if isObjectNotFoundError(err) {
			return r.diffByLinearHistory(ctx, branchTip, baseTip)
		}
		return nil, fmt.Errorf("iterate branch commits: %w", err)
	}

	return out, nil
}

func (r *Repository) diffByLinearHistory(ctx context.Context, branchTip, baseTip plumbing.Hash) ([]Commit, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	iter, err := r.repo.Log(&git.LogOptions{From: branchTip})
	if err != nil {
		return nil, fmt.Errorf("list branch commits: %w", err)
	}
	defer iter.Close()

	var out []Commit
	reachedBase := false
	err = iter.ForEach(func(c *object.Commit) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if c.Hash == baseTip {
			reachedBase = true
			return storer.ErrStop
		}
		out = append(out, toCommit(c))
		return nil
	})
	if err != nil && !errors.Is(err, storer.ErrStop) {
		if isObjectNotFoundError(err) && len(out) > 0 {
			return out, nil
		}
		return nil, fmt.Errorf("iterate branch commits: %w", err)
	}
	if !reachedBase && len(out) > 0 {
		return out, nil
	}

	return out, nil
}

func (r *Repository) commitsFrom(ctx context.Context, from plumbing.Hash) ([]Commit, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	iter, err := r.repo.Log(&git.LogOptions{From: from})
	if err != nil {
		return nil, fmt.Errorf("list commits from %s: %w", from.String(), err)
	}
	defer iter.Close()

	var out []Commit
	err = iter.ForEach(func(c *object.Commit) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		out = append(out, toCommit(c))
		return nil
	})
	if err != nil {
		if isObjectNotFoundError(err) && len(out) > 0 {
			return out, nil
		}
		return nil, fmt.Errorf("iterate commits from %s: %w", from.String(), err)
	}

	return out, nil
}

func (r *Repository) reachableSet(ctx context.Context, from plumbing.Hash) (map[plumbing.Hash]struct{}, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	iter, err := r.repo.Log(&git.LogOptions{From: from})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	out := map[plumbing.Hash]struct{}{}
	err = iter.ForEach(func(c *object.Commit) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		out[c.Hash] = struct{}{}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (r *Repository) resolveRevision(value string) (plumbing.Hash, error) {
	rev := plumbing.Revision(strings.TrimSpace(value))
	resolved, err := r.repo.ResolveRevision(rev)
	if err != nil {
		return plumbing.ZeroHash, err
	}
	return *resolved, nil
}

func (r *Repository) resolveBranchTip(branch string) (plumbing.Hash, error) {
	if strings.TrimSpace(branch) == "" {
		return plumbing.ZeroHash, fmt.Errorf("branch must not be empty")
	}

	candidates := []string{
		branch,
		"refs/heads/" + branch,
		"refs/remotes/" + branch,
		"origin/" + branch,
		"refs/remotes/origin/" + branch,
	}

	for _, candidate := range candidates {
		h, err := r.resolveRevision(candidate)
		if err == nil {
			return h, nil
		}
	}

	return plumbing.ZeroHash, fmt.Errorf("branch %q not found", branch)
}

func toCommit(c *object.Commit) Commit {
	subject := c.Message
	if idx := strings.IndexByte(subject, '\n'); idx >= 0 {
		subject = subject[:idx]
	}
	full := c.Hash.String()
	short := full
	if len(short) > 7 {
		short = short[:7]
	}

	return Commit{
		Hash:      full,
		ShortHash: short,
		Subject:   strings.TrimSpace(subject),
		Message:   strings.TrimRight(c.Message, "\r\n"),
	}
}

func isObjectNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "object not found")
}
