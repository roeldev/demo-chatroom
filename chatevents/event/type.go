// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package event

import (
	"github.com/roeldev/demo-chatroom/chatusers"
)

// Type restricts which structs can be used within a [chatevents.Event].
type Type interface {
	eventType()
}

// A UserEvent provides user related data.
type UserEvent interface {
	Type
	GetUserID() chatusers.UserID
	GetUserDetails() chatusers.UserDetails
}

// A ReceiverEvent provides receiver related data.
type ReceiverEvent interface {
	Type
	GetReceiverID() chatusers.UserID
}

var _ Type = (*event)(nil)

// event is used to mark a struct as a [Type].
type event struct{}

func (*event) eventType() {}
