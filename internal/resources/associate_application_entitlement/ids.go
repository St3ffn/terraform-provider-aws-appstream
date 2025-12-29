// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_application_entitlement

import (
	"fmt"
	"strings"
)

func buildID(stackName, entitlementName, applicationIdentifier string) string {
	return fmt.Sprintf("%s|%s|%s", stackName, entitlementName, applicationIdentifier)
}

func parseID(id string) (stackName, entitlementName, applicationIdentifier string, err error) {
	parts := strings.SplitN(id, "|", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("invalid associate application entitlement ID format")
	}

	return parts[0], parts[1], parts[2], nil
}
