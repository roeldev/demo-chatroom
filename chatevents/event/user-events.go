// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package event

import (
	"github.com/roeldev/demo-chatroom/chatusers"
)

var (
	_ UserEvent = (*UserJoinEvent)(nil)
	_ UserEvent = (*UserLeaveEvent)(nil)
	_ UserEvent = (*UserUpdateEvent)(nil)
	_ UserEvent = (*UserStatusEvent)(nil)
	_ UserEvent = (*UserTypingEvent)(nil)
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

type UserTypingEvent struct {
	event
	UserID      chatusers.UserID
	UserDetails chatusers.UserDetails
	ReceiverID  chatusers.UserID
	IsTyping    bool
}

func (e UserJoinEvent) GetUserID() chatusers.UserID   { return e.UserID }
func (e UserLeaveEvent) GetUserID() chatusers.UserID  { return e.UserID }
func (e UserUpdateEvent) GetUserID() chatusers.UserID { return e.UserID }
func (e UserStatusEvent) GetUserID() chatusers.UserID { return e.UserID }
func (e UserTypingEvent) GetUserID() chatusers.UserID { return e.UserID }

func (e UserJoinEvent) GetUserDetails() chatusers.UserDetails   { return e.UserDetails }
func (e UserLeaveEvent) GetUserDetails() chatusers.UserDetails  { return e.UserDetails }
func (e UserUpdateEvent) GetUserDetails() chatusers.UserDetails { return e.After }
func (e UserStatusEvent) GetUserDetails() chatusers.UserDetails { return e.UserDetails }
func (e UserTypingEvent) GetUserDetails() chatusers.UserDetails { return e.UserDetails }
