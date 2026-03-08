// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/st3ffn/terraform-provider-aws-appstream/tools/commitlint/internal/model"
)

func TestWriteTextAndJSONAndGitHubAndMarkdown(t *testing.T) {
	t.Parallel()

	report := model.Report{
		TotalCommits:   2,
		CheckedCommits: 2,
		InvalidCommits: 1,
		Results: []model.CommitResult{
			{Hash: "aaaaaaaa", ShortHash: "aaaaaaa", Subject: "feat: ok", Message: "feat: ok", Valid: true},
			{Hash: "bbbbbbbb", ShortHash: "bbbbbbb", Subject: "bad", Message: "bad", Valid: false, Violations: []model.Violation{{Code: "invalid-header", Message: "header invalid"}}},
		},
	}

	for _, format := range []Format{FormatText, FormatJSON, FormatGitHub, FormatMarkdown, FormatRDJSONL} {
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			var b bytes.Buffer
			err := Write(report, Options{Format: format, ColorMode: ColorNever, Writer: &b})
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			out := b.String()
			if strings.TrimSpace(out) == "" {
				t.Fatal("expected non-empty output")
			}
		})
	}
}

func TestWriteRDJSONL(t *testing.T) {
	t.Parallel()

	report := model.Report{
		Results: []model.CommitResult{
			{
				ShortHash: "abc1234",
				Subject:   "bad",
				Message:   "bad",
				Valid:     false,
				Violations: []model.Violation{
					{Code: "type-empty", Message: "type may not be empty"},
					{Code: "subject-empty", Message: "subject may not be empty"},
				},
			},
		},
	}

	var b bytes.Buffer
	if err := Write(report, Options{Format: FormatRDJSONL, Writer: &b}); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(b.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}

	for _, line := range lines {
		var diag map[string]any
		if err := json.Unmarshal([]byte(line), &diag); err != nil {
			t.Fatalf("invalid json line %q: %v", line, err)
		}
		if diag["severity"] != "ERROR" {
			t.Fatalf("severity=%v want ERROR", diag["severity"])
		}
	}
}

func TestWriteUnsupportedFormat(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	err := Write(model.Report{}, Options{Format: Format("unknown"), Writer: &b})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
