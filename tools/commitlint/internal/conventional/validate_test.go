// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package conventional

import "testing"

func TestValidateCommitMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		msg     string
		wantErr bool
	}{
		{name: "valid_simple", msg: "feat: add resource", wantErr: false},
		{name: "valid_scope", msg: "fix(parser): handle empty body", wantErr: false},
		{name: "valid_breaking_bang", msg: "feat(api)!: remove old endpoint", wantErr: false},
		{name: "valid_body", msg: "docs: update readme\n\nadd details", wantErr: false},
		{name: "valid_breaking_footer", msg: "refactor: switch default\n\nBREAKING CHANGE: old field removed", wantErr: false},
		{name: "valid_initial_commit", msg: "Initial commit", wantErr: false},
		{name: "invalid_header", msg: "not conventional", wantErr: true},
		{name: "invalid_no_blank_line", msg: "feat: title\nbody without blank line", wantErr: true},
		{name: "invalid_breaking_footer_empty", msg: "feat: title\n\nBREAKING CHANGE:", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			errs := ValidateCommitMessage(tt.msg)
			gotErr := len(errs) > 0
			if gotErr != tt.wantErr {
				t.Fatalf("ValidateCommitMessage() error=%v, wantErr=%v (errs=%v)", gotErr, tt.wantErr, errs)
			}
		})
	}
}

func TestIsMergeCommitSubject(t *testing.T) {
	t.Parallel()

	if !IsMergeCommitSubject("Merge branch 'main' into feature") {
		t.Fatal("expected merge subject to be detected")
	}
	if IsMergeCommitSubject("feat: add command") {
		t.Fatal("did not expect conventional commit subject to be merge")
	}
}
