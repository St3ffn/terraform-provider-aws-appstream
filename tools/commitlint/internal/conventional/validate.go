// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package conventional

import (
	"regexp"
	"strings"

	"github.com/st3ffn/terraform-provider-aws-appstream/tools/commitlint/internal/model"
)

var headerRegexp = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9-]*)(\([^)\r\n]+\))?(!)?: ([^\r\n]+)$`)

// IsMergeCommitSubject reports whether a commit subject is a merge commit title.
func IsMergeCommitSubject(subject string) bool {
	s := strings.TrimSpace(subject)
	return strings.HasPrefix(s, "Merge ")
}

// ValidateCommitMessage validates a full commit message against Conventional Commits.
//
// Rules enforced:
//   - Header must match <type>[optional scope][!]: <description>
//   - Exception: "Initial commit" is accepted as a valid header
//   - If additional lines exist, line 2 must be blank
//   - BREAKING CHANGE footer (when present) must contain text after colon
func ValidateCommitMessage(message string) []model.Violation {
	var violations []model.Violation

	trimmed := strings.TrimRight(message, "\r\n")
	if strings.TrimSpace(trimmed) == "" {
		return []model.Violation{{
			Code:    "empty-message",
			Message: "commit message must not be empty",
		}}
	}

	lines := splitLines(trimmed)
	header := lines[0]
	if !isInitialCommitHeader(header) && !headerRegexp.MatchString(header) {
		violations = append(violations, model.Violation{
			Code:    "invalid-header",
			Message: "header must match <type>[optional scope][!]: <description>",
		})
	}

	if len(lines) > 1 && strings.TrimSpace(lines[1]) != "" {
		violations = append(violations, model.Violation{
			Code:    "missing-blank-line-after-header",
			Message: "body/footer must start after a blank line",
		})
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "BREAKING CHANGE:") {
			if strings.TrimSpace(strings.TrimPrefix(line, "BREAKING CHANGE:")) == "" {
				violations = append(violations, model.Violation{
					Code:    "invalid-breaking-change-footer",
					Message: "BREAKING CHANGE footer must include description text",
				})
			}
		}
	}

	return violations
}

func isInitialCommitHeader(header string) bool {
	return strings.TrimSpace(header) == "Initial commit"
}

func splitLines(input string) []string {
	normalized := strings.ReplaceAll(input, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	return strings.Split(normalized, "\n")
}
