// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package model

// Violation represents a single Conventional Commit rule violation.
type Violation struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CommitResult captures lint results for one commit.
type CommitResult struct {
	Hash       string      `json:"hash"`
	ShortHash  string      `json:"short_hash"`
	Subject    string      `json:"subject"`
	Message    string      `json:"message"`
	Valid      bool        `json:"valid"`
	Violations []Violation `json:"violations,omitempty"`
}

// Report summarizes a lint run and contains per-commit results.
type Report struct {
	TotalCommits   int            `json:"total_commits"`
	CheckedCommits int            `json:"checked_commits"`
	SkippedCommits int            `json:"skipped_commits"`
	InvalidCommits int            `json:"invalid_commits"`
	Results        []CommitResult `json:"results"`
}
