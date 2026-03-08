// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package run

import "testing"

func TestValidateOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    Options
		wantErr bool
	}{
		{name: "default_ok", opts: Options{}, wantErr: false},
		{name: "all_ok", opts: Options{All: true}, wantErr: false},
		{name: "branch_ok", opts: Options{Branch: "feature"}, wantErr: false},
		{name: "branch_diff_ok", opts: Options{BranchDiff: "feature", Base: "main"}, wantErr: false},
		{name: "range_ok", opts: Options{From: "a", To: "b"}, wantErr: false},
		{name: "range_missing_from", opts: Options{To: "b"}, wantErr: true},
		{name: "range_missing_to", opts: Options{From: "a"}, wantErr: true},
		{name: "mutually_exclusive_all_branch", opts: Options{All: true, Branch: "feature"}, wantErr: true},
		{name: "mutually_exclusive_branch_range", opts: Options{Branch: "feature", From: "a", To: "b"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateOptions(tt.opts)
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Fatalf("validateOptions() error=%v wantErr=%v err=%v", gotErr, tt.wantErr, err)
			}
		})
	}
}

func TestApplyDefaults(t *testing.T) {
	t.Parallel()

	opts := applyDefaults(Options{})
	if opts.RepoPath != defaultRepoPath {
		t.Fatalf("RepoPath=%q want %q", opts.RepoPath, defaultRepoPath)
	}
	if opts.Base != defaultBaseRef {
		t.Fatalf("Base=%q want %q", opts.Base, defaultBaseRef)
	}
}
