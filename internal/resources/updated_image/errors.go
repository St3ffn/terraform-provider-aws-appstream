// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package updated_image

import "errors"

var errUnexpectedUpdatedImageState = errors.New("unexpected updated image state")

func isUnexpectedUpdatedImageStateError(err error) bool {
	return errors.Is(err, errUnexpectedUpdatedImageState)
}
