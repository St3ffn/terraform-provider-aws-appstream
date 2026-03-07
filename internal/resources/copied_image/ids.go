// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package copied_image

import (
	"fmt"
	"strings"
)

func buildID(destinationImageName, destinationRegion, sourceImageName string) string {
	return fmt.Sprintf("%s|%s|%s", destinationImageName, destinationRegion, sourceImageName)
}

func parseID(id string) (destinationImageName, destinationRegion, sourceImageName string, err error) {
	parts := strings.Split(id, "|")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("invalid copied image import ID format")
	}

	return parts[0], parts[1], parts[2], nil
}
