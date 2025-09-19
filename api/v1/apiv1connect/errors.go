// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1connect

import "github.com/go-pogo/errors"

const (
	ErrInvalidToken      errors.Msg = "invalid token"
	ErrInvalidUserID     errors.Msg = "invalid user id"
	ErrInvalidChatID     errors.Msg = "invalid chat id"
	ErrInvalidReceiverID errors.Msg = "invalid receiver id"
	ErrChangeUserStatus  errors.Msg = "failed to change user status"
)
