// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1

import (
	"github.com/roeldev/demo-chatroom/chatevents/event"
	"github.com/roeldev/demo-chatroom/chatusers"
)

type EventStreamResponseEvent = isEventStreamResponse_Event

func NewEventUser(u event.UserEvent) *EventUser {
	return &EventUser{
		Id:      NewUUID(u.GetUserID()),
		Details: NewUserDetails(u.GetUserDetails()),
	}
}

func NewUserMentions(users map[chatusers.UserID]string) []*UserMention {
	mentions := make([]*UserMention, 0, len(users))
	for uid, name := range users {
		mentions = append(mentions, &UserMention{
			UserId:   NewUUID(uid),
			UserName: name,
		})
	}
	return mentions
}
