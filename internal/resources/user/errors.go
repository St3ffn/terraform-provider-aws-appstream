// Copyright (c) St3ffn
// SPDX-License-Identifier: MPL-2.0

package user

import "errors"

var errUserNotYetVisible = errors.New("appstream user not yet visible")

func isUserNotYetVisibleError(err error) bool {
	return errors.Is(err, errUserNotYetVisible)
}
