// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"strings"
)

func buildEntitlementID(stackName, name string) string {
	return fmt.Sprintf("%s|%s", stackName, name)
}

func parseEntitlementID(id string) (stackName, name string, err error) {
	parts := strings.SplitN(id, "|", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid entitlement ID format")
	}

	return parts[0], parts[1], nil
}
