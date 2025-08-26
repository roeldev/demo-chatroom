// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1

import (
	"strconv"

	"github.com/roeldev/demo-chatroom/chatusers"
)

func FromJoinRequest(x *JoinRequest) chatusers.UserOption {
	return func(u *chatusers.User) error {
		u.Name = x.User.Name
		u.Initials = x.User.Initials

		var err error
		if u.Color1, err = x.User.Color1.Decode(); err != nil {
			return err
		}
		if u.Color2, err = x.User.Color2.Decode(); err != nil {
			return err
		}

		u.Flags = x.Flags.ToChatUserFlags()
		return nil
	}
}

func NewUserDetails(u chatusers.UserDetails) *UserDetails {
	return &UserDetails{
		Name:     u.Name,
		Initials: u.Initials,
		Color1:   NewColor(u.Color1),
		Color2:   NewColor(u.Color2),
	}
}

func NewUserFlags(flag chatusers.Flag) UserFlag {
	switch flag {
	case chatusers.Flag_None:
		return UserFlag_USER_FLAG_NONE
	case chatusers.Flag_IsBot:
		return UserFlag_USER_FLAG_IS_BOT
	default:
		panic("apiv1: Flag " + itoa(int64(flag)) + " is invalid")
	}
}

func (x UserFlag) ToChatUserFlags() chatusers.Flag {
	switch x {
	case UserFlag_USER_FLAG_NONE:
		return chatusers.Flag_None
	case UserFlag_USER_FLAG_IS_BOT:
		return chatusers.Flag_IsBot
	default:
		panic("apiv1: UserFlag " + itoa(int64(x)) + " is invalid")
	}
}

func NewUserStatus(stat chatusers.Status) UserStatus {
	switch stat {
	case chatusers.Status_Default:
		return UserStatus_USER_STATUS_DEFAULT
	case chatusers.Status_Busy:
		return UserStatus_USER_STATUS_BUSY
	case chatusers.Status_Away:
		return UserStatus_USER_STATUS_AWAY
	default:
		panic("apiv1: Status " + itoa(int64(stat)) + " is invalid")
	}
}

func (x UserStatus) ToChatUserStatus() chatusers.Status {
	switch x {
	case UserStatus_USER_STATUS_DEFAULT:
		return chatusers.Status_Default
	case UserStatus_USER_STATUS_BUSY:
		return chatusers.Status_Busy
	case UserStatus_USER_STATUS_AWAY:
		return chatusers.Status_Away
	default:
		panic("apiv1: UserStatus " + itoa(int64(x)) + " is invalid")
	}
}

func itoa(v int64) string { return strconv.FormatInt(v, 10) }
