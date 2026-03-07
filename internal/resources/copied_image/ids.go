// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package copied_image

import (
	"fmt"
	"strings"
)

func buildID(name, destinationRegion string) string {
	return fmt.Sprintf("%s|%s", name, destinationRegion)
}

func parseID(id string) (name, destinationRegion string, err error) {
	parts := strings.SplitN(id, "|", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid copied image ID format")
	}

	return parts[0], parts[1], nil
}
