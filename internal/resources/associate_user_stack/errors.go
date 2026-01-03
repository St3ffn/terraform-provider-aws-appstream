// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package associate_user_stack

import (
	"errors"
	"fmt"

	awstypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
)

var errUserStackAssociationNotReady = errors.New("appstream user-stack association not yet ready")

func newUserStackAssociationNotReadyError(e awstypes.UserStackAssociationError) error {
	msg := "unknown appstream user-stack association error"
	if e.ErrorMessage != nil {
		msg = *e.ErrorMessage
	}

	return fmt.Errorf("%w: %s", errUserStackAssociationNotReady, msg)
}

func isUserStackAssociationNotReadyError(err error) bool {
	return errors.Is(err, errUserStackAssociationNotReady)
}
