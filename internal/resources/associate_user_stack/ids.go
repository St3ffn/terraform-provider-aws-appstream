// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack

import (
	"fmt"
	"strings"
)

func buildID(stackName, authenticationType, userName string) string {
	return fmt.Sprintf("%s|%s|%s", stackName, authenticationType, userName)
}

func parseID(id string) (stackName, authenticationType, userName string, err error) {
	parts := strings.SplitN(id, "|", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("invalid associate user stack ID format")
	}

	return parts[0], parts[1], parts[2], nil
}
