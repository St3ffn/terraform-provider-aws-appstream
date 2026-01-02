// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_application_fleet

import (
	"fmt"
	"strings"
)

func buildID(fleetName, applicationARN string) string {
	return fmt.Sprintf("%s|%s", fleetName, applicationARN)
}

func parseID(id string) (fleetName, applicationARN string, err error) {
	parts := strings.SplitN(id, "|", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid associate application fleet ID format")
	}

	return parts[0], parts[1], nil
}
