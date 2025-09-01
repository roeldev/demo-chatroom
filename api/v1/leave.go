// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1

import "github.com/roeldev/demo-chatroom/chatevents/event"

func NewLeaveReason(reason event.LeaveReason) LeaveReason {
	switch reason {
	case event.UserLeave:
		return LeaveReason_LEAVE_REASON_USER_ACTION

	case event.Disconnected:
		return LeaveReason_LEAVE_REASON_DISCONNECTED

	default:
		panic("apiv1: LeaveReason " + itoa(int64(reason)) + " is invalid")
	}
}
