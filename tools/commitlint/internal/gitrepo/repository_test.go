// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package gitrepo

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestSelectBranchDiffAndBranch(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	r, hashes := seedRepository(t, dir)

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	diff, err := repo.SelectBranchDiff(context.Background(), "feature", "main")
	if err != nil {
		t.Fatalf("SelectBranchDiff() error = %v", err)
	}
	if len(diff) != 2 {
		t.Fatalf("SelectBranchDiff() len = %d, want 2", len(diff))
	}
	if diff[0].Hash != hashes["feature2"].String() || diff[1].Hash != hashes["feature1"].String() {
		t.Fatalf("unexpected branch diff order/hashes: %#v", diff)
	}

	branchAll, err := repo.SelectBranch(context.Background(), "feature")
	if err != nil {
		t.Fatalf("SelectBranch() error = %v", err)
	}
	if len(branchAll) < 4 {
		t.Fatalf("SelectBranch() len = %d, want at least 4", len(branchAll))
	}

	if err := checkoutBranch(r, "feature"); err != nil {
		t.Fatalf("checkout feature failed: %v", err)
	}
	currentDiff, err := repo.SelectCurrentBranchDiff(context.Background(), "main")
	if err != nil {
		t.Fatalf("SelectCurrentBranchDiff() error = %v", err)
	}
	if len(currentDiff) != 2 {
		t.Fatalf("SelectCurrentBranchDiff() len = %d, want 2", len(currentDiff))
	}
}

func TestSelectRangeInclusive(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	_, hashes := seedRepository(t, dir)

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	commits, err := repo.SelectRange(context.Background(), hashes["feature1"].String(), hashes["feature2"].String())
	if err != nil {
		t.Fatalf("SelectRange() error = %v", err)
	}
	if len(commits) != 2 {
		t.Fatalf("SelectRange() len = %d, want 2", len(commits))
	}
	if commits[0].Hash != hashes["feature2"].String() || commits[1].Hash != hashes["feature1"].String() {
		t.Fatalf("SelectRange() unexpected commits: %#v", commits)
	}
}

func TestSelectAll(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	_, _ = seedRepository(t, dir)

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	commits, err := repo.SelectAll(context.Background())
	if err != nil {
		t.Fatalf("SelectAll() error = %v", err)
	}
	if len(commits) < 4 {
		t.Fatalf("SelectAll() len = %d, want at least 4", len(commits))
	}
}

func TestSelectCurrentBranchDiffWhenHeadEqualsBase(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	_, _ = seedRepository(t, dir)

	repo, err := Open(dir)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	commits, err := repo.SelectCurrentBranchDiff(context.Background(), "main")
	if err != nil {
		t.Fatalf("SelectCurrentBranchDiff() error = %v", err)
	}
	if len(commits) != 0 {
		t.Fatalf("SelectCurrentBranchDiff() len = %d, want 0", len(commits))
	}
}

func seedRepository(t *testing.T, dir string) (*git.Repository, map[string]plumbing.Hash) {
	t.Helper()

	r, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("PlainInit() error = %v", err)
	}

	wt, err := r.Worktree()
	if err != nil {
		t.Fatalf("Worktree() error = %v", err)
	}

	hashes := map[string]plumbing.Hash{}
	hashes["main1"] = commitFile(t, wt, dir, "a.txt", "a1", "feat: first main commit")
	hashes["main2"] = commitFile(t, wt, dir, "a.txt", "a2", "fix: second main commit")

	if err := wt.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName("main"), Create: true}); err != nil {
		t.Fatalf("create main branch failed: %v", err)
	}

	if err := wt.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName("feature"), Create: true}); err != nil {
		t.Fatalf("checkout feature failed: %v", err)
	}

	hashes["feature1"] = commitFile(t, wt, dir, "b.txt", "b1", "feat: feature one")
	hashes["feature2"] = commitFile(t, wt, dir, "b.txt", "b2", "feat: feature two")

	if err := checkoutBranch(r, "main"); err != nil {
		t.Fatalf("checkout main failed: %v", err)
	}

	return r, hashes
}

func checkoutBranch(r *git.Repository, branch string) error {
	wt, err := r.Worktree()
	if err != nil {
		return err
	}
	return wt.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(branch)})
}

func commitFile(t *testing.T, wt *git.Worktree, dir, file, content, msg string) plumbing.Hash {
	t.Helper()

	path := filepath.Join(dir, file)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", file, err)
	}
	if _, err := wt.Add(file); err != nil {
		t.Fatalf("Add(%s) error = %v", file, err)
	}

	hash, err := wt.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "commitlint-test",
			Email: "commitlint@example.invalid",
			When:  time.Now().UTC(),
		},
	})
	if err != nil {
		t.Fatalf("Commit(%q) error = %v", msg, err)
	}

	return hash
}
