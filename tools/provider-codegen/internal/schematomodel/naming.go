// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package schematomodel

import (
	"strings"
	"unicode"
)

func toGoFieldName(tfName string) string {
	parts := strings.Split(tfName, "_")
	for i, part := range parts {
		parts[i] = normalizeWord(part)
	}

	return strings.Join(parts, "")
}

func nestedModelTypeName(rootTypeName string, path []string) string {
	if len(path) == 0 {
		return rootTypeName
	}

	parts := make([]string, 0, len(path))
	for _, part := range path {
		parts = append(parts, toGoFieldName(part))
	}

	s := strings.Join(parts, "")
	if s == "" {
		return rootTypeName
	}

	return rootTypeName + s
}

func lowerFirst(s string) string {
	rs := []rune(s)
	if len(rs) == 0 {
		return s
	}

	// Keep leading acronyms in lower camel case: VPCConfig -> vpcConfig.
	last := 1
	for last < len(rs) && unicode.IsUpper(rs[last]) {
		last++
		if last < len(rs) && unicode.IsLower(rs[last]) {
			last--
			break
		}
	}

	for i := 0; i < last; i++ {
		rs[i] = unicode.ToLower(rs[i])
	}

	return string(rs)
}

func upperFirst(s string) string {
	rs := []rune(s)
	if len(rs) == 0 {
		return s
	}
	rs[0] = unicode.ToUpper(rs[0])
	return string(rs)
}

func normalizeWord(word string) string {
	if word == "" {
		return ""
	}

	acronyms := map[string]string{
		"arn":  "ARN",
		"id":   "ID",
		"ids":  "IDs",
		"iam":  "IAM",
		"imds": "IMDS",
		"gb":   "GB",
		"s3":   "S3",
		"usb":  "USB",
		"vpc":  "VPC",
	}

	lower := strings.ToLower(word)
	if acronym, ok := acronyms[lower]; ok {
		return acronym
	}

	rs := []rune(lower)
	rs[0] = unicode.ToUpper(rs[0])
	return string(rs)
}
