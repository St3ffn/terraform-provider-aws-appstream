// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package image_builder

import "errors"

var errUnexpectedImageBuilderState = errors.New("unexpected image builder state")

func isUnexpectedImageBuilderStateError(err error) bool {
	return errors.Is(err, errUnexpectedImageBuilderState)
}
