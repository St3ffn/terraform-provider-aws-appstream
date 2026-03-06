// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package imported_image

import "errors"

var errUnexpectedImportedImageState = errors.New("unexpected imported image state")

func isUnexpectedImportedImageStateError(err error) bool {
	return errors.Is(err, errUnexpectedImportedImageState)
}
