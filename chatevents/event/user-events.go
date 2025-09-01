// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package event

import (
	"github.com/google/uuid"
	"github.com/roeldev/demo-chatroom/chatusers"
	"github.com/rs/zerolog"
)

type LeaveReason uint8

//goland:noinspection GoSnakeCaseUsage
const (
	UserLeave LeaveReason = iota
	Disconnected
)

var (
	_ UserEvent = (*UserJoinEvent)(nil)
	_ UserEvent = (*UserLeaveEvent)(nil)
	_ UserEvent = (*UserUpdateEvent)(nil)
	_ UserEvent = (*UserStatusEvent)(nil)
)

type UserJoinEvent struct {
	event
	UserID      chatusers.UserID
	UserDetails chatusers.UserDetails
	UserFlags   chatusers.Flag
}

type UserLeaveEvent struct {
	event
	UserID      chatusers.UserID
	UserDetails chatusers.UserDetails
	Reason      LeaveReason
}

type UserUpdateEvent struct {
	event
	UserID chatusers.UserID
	Before chatusers.UserDetails
	After  chatusers.UserDetails
}

type UserStatusEvent struct {
	event
	UserID      chatusers.UserID
	UserDetails chatusers.UserDetails
	Before      chatusers.Status
	After       chatusers.Status
}

func (e UserJoinEvent) GetUserID() chatusers.UserID   { return e.UserID }
func (e UserLeaveEvent) GetUserID() chatusers.UserID  { return e.UserID }
func (e UserUpdateEvent) GetUserID() chatusers.UserID { return e.UserID }
func (e UserStatusEvent) GetUserID() chatusers.UserID { return e.UserID }

func (e UserJoinEvent) GetUserDetails() chatusers.UserDetails   { return e.UserDetails }
func (e UserLeaveEvent) GetUserDetails() chatusers.UserDetails  { return e.UserDetails }
func (e UserUpdateEvent) GetUserDetails() chatusers.UserDetails { return e.After }
func (e UserStatusEvent) GetUserDetails() chatusers.UserDetails { return e.UserDetails }

var (
	_ zerolog.LogObjectMarshaler = (*UserTypingEvent)(nil)
	_ UserEvent                  = (*UserTypingEvent)(nil)
	_ ReceiverEvent              = (*UserTypingEvent)(nil)
)

type UserTypingEvent struct {
	event
	UserID      chatusers.UserID
	UserDetails chatusers.UserDetails
	ReceiverID  chatusers.UserID
	IsTyping    bool
}

func (e UserTypingEvent) GetUserID() chatusers.UserID           { return e.UserID }
func (e UserTypingEvent) GetUserDetails() chatusers.UserDetails { return e.UserDetails }
func (e UserTypingEvent) GetReceiverID() chatusers.UserID       { return e.ReceiverID }

func (e UserTypingEvent) MarshalZerologObject(ze *zerolog.Event) {
	ze.Str("user", chatusers.IdentifierString(e.UserID, e.UserDetails))
	if e.ReceiverID != uuid.Nil {
		ze.Stringer("receiver_id", e.ReceiverID)
	}

	ze.Bool("typing", e.IsTyping)
}
