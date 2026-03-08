// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package imported_image

import (
	"fmt"
	"strings"
)

func buildID(name, iamRoleARN, sourceAmiID string) string {
	return fmt.Sprintf("%s|%s|%s", name, iamRoleARN, sourceAmiID)
}

func parseID(id string) (name, iamRoleARN, sourceAmiID string, err error) {
	parts := strings.Split(id, "|")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("invalid imported image import ID format")
	}

	return parts[0], parts[1], parts[2], nil
}
