// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package run

import (
	"errors"
)

const (
	defaultRepoPath = "."
	defaultBaseRef  = "main"
)

// Options configures commit selection and lint behavior.
type Options struct {
	RepoPath            string
	All                 bool
	Branch              string
	BranchDiff          string
	Base                string
	From                string
	To                  string
	IncludeMergeCommits bool
}

func applyDefaults(opts Options) Options {
	if opts.RepoPath == "" {
		opts.RepoPath = defaultRepoPath
	}
	if opts.Base == "" {
		opts.Base = defaultBaseRef
	}
	return opts
}

func validateOptions(opts Options) error {
	if (opts.From == "") != (opts.To == "") {
		return errors.New("--from and --to must be set together")
	}

	selectors := 0
	if opts.All {
		selectors++
	}
	if opts.Branch != "" {
		selectors++
	}
	if opts.BranchDiff != "" {
		selectors++
	}
	if opts.From != "" && opts.To != "" {
		selectors++
	}

	if selectors > 1 {
		return errors.New("commit selection flags are mutually exclusive: use only one of --all, --branch, --branch-diff, or --from/--to")
	}

	return nil
}
