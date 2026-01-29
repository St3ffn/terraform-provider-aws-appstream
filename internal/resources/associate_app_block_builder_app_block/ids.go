// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package associate_app_block_builder_app_block

import (
	"fmt"
	"strings"
)

func buildID(appBlockBuilderName, appBlockARN string) string {
	return fmt.Sprintf("%s|%s", appBlockBuilderName, appBlockARN)
}

func parseID(id string) (appBlockBuilderName, appBlockARN string, err error) {
	parts := strings.SplitN(id, "|", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid associate app block builder app block ID format")
	}

	return parts[0], parts[1], nil
}
