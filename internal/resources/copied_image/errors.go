// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package copied_image

import "errors"

var errUnexpectedCopiedImageState = errors.New("unexpected copied image state")

func isUnexpectedCopiedImageStateError(err error) bool {
	return errors.Is(err, errUnexpectedCopiedImageState)
}
