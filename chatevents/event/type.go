// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package event

import "github.com/roeldev/demo-chatroom/chatusers"

type Type interface {
	eventType()
}

var _ Type = (*event)(nil)

type event struct{}

func (*event) eventType() {}

type UserEvent interface {
	Type
	GetUserID() chatusers.UserID
	GetUserDetails() chatusers.UserDetails
}
