// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package fleet

import "errors"

var errUnexpectedFleetState = errors.New("unexpected fleet state")

func isUnexpectedFleetStateError(err error) bool {
	return errors.Is(err, errUnexpectedFleetState)
}
