// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package updated_image

import (
	"fmt"
	"strings"
)

func buildID(existingImageName, newImageName string) string {
	return fmt.Sprintf("%s|%s", existingImageName, newImageName)
}

func parseID(id string) (existingImageName, newImageName string, err error) {
	parts := strings.Split(id, "|")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid updated image import ID format")
	}

	return parts[0], parts[1], nil
}
