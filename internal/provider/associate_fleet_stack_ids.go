// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"strings"
)

func buildAssociateFleetStackID(fleetName, stackName string) string {
	return fmt.Sprintf("%s|%s", fleetName, stackName)
}

func parseAssociateFleetStackID(id string) (fleetName, stackName string, err error) {
	parts := strings.SplitN(id, "|", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid associate fleet stack ID format")
	}

	return parts[0], parts[1], nil
}
