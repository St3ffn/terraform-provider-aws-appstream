// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import (
	"go/ast"
	"regexp"
	"strconv"
	"strings"
)

var codegenHasRemoteChangesAnnotationRE = regexp.MustCompile(
	`\bcodegen:has_remote_changes\s*=\s*(true|false)\b`,
)

func remoteOverrideFromCommentMap(cm ast.CommentMap, kv *ast.KeyValueExpr) *bool {
	var (
		found bool
		value bool
	)

	scanGroups := func(groups []*ast.CommentGroup) {
		for _, group := range groups {
			v, ok := parseRemoteOverride(group.Text())
			if !ok {
				continue
			}
			found = true
			value = v
		}
	}

	scanGroups(cm[kv])
	scanGroups(cm[kv.Key])
	scanGroups(cm[kv.Value])

	if !found {
		return nil
	}

	return &value
}

func parseRemoteOverride(text string) (bool, bool) {
	for _, line := range strings.Split(text, "\n") {
		match := codegenHasRemoteChangesAnnotationRE.FindStringSubmatch(line)
		if len(match) != 2 {
			continue
		}

		v, err := strconv.ParseBool(match[1])
		if err != nil {
			return false, false
		}

		return v, true
	}

	return false, false
}
