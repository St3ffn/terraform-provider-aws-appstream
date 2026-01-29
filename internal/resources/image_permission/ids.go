// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package image_permission

import (
	"fmt"
	"strings"
)

func buildID(name, sharedAccountID string) string {
	return fmt.Sprintf("%s|%s", name, sharedAccountID)
}

func parseID(id string) (name, sharedAccountID string, err error) {
	parts := strings.SplitN(id, "|", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid image permission ID format")
	}

	return parts[0], parts[1], nil
}
