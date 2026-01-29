// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package app_block_builder

import "errors"

var errUnexpectedAppBlockBuilderState = errors.New("unexpected app block builder state")

func isUnexpectedAppBlockBuilderStateError(err error) bool {
	return errors.Is(err, errUnexpectedAppBlockBuilderState)
}
