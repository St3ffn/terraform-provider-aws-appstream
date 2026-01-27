// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_image_builder_software

import "errors"

var errUnexpectedImageBuilderState = errors.New("unexpected image builder state")

func isUnexpectedImageBuilderStateError(err error) bool {
	return errors.Is(err, errUnexpectedImageBuilderState)
}
