// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"fmt"
	"strings"
)

func buildID(authenticationType, userName string) string {
	return fmt.Sprintf("%s|%s", authenticationType, userName)
}

func parseID(id string) (authenticationType, userName string, err error) {
	parts := strings.SplitN(id, "|", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid user ID format")
	}

	return parts[0], parts[1], nil
}
