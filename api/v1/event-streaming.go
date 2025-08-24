// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1

import (
	"reflect"

	"github.com/roeldev/demo-chatroom/chatevents/event"
)

func NewEventStreamResponseEvent(typ event.Type) EventStreamResponseEvent {
	switch et := typ.(type) {
	case *event.UserJoinEvent:
		return &EventStreamResponse_UserJoin{
			UserJoin: &UserJoinEvent{
				User:  NewEventUser(et),
				Flags: NewUserFlags(et.UserFlags),
			},
		}

	case *event.UserLeaveEvent:
		return &EventStreamResponse_UserLeave{
			UserLeave: &UserLeaveEvent{
				User: NewEventUser(et),
			},
		}

	case *event.UserUpdateEvent:
		return &EventStreamResponse_UserUpdate{
			UserUpdate: &UserUpdateEvent{
				User:   NewEventUser(et),
				Before: NewUserDetails(et.Before),
			},
		}

	case *event.UserStatusEvent:
		return &EventStreamResponse_UserStatus{
			UserStatus: &UserStatusEvent{
				User:   NewEventUser(et),
				Before: NewUserStatus(et.Before),
			},
		}

	case *event.UserTypingEvent:
		return &EventStreamResponse_UserTyping{
			UserTyping: &UserTypingEvent{
				User:       NewEventUser(et),
				ReceiverId: NewUUID(et.ReceiverID),
				Typing:     et.IsTyping,
			},
		}

	case *event.ChatEvent:
		return &EventStreamResponse_ChatSent{
			ChatSent: &ChatSentEvent{
				ChatId:      NewUUID(et.ChatID),
				User:        NewEventUser(et),
				ReceiverId:  NewUUID(et.ReceiverID),
				ReplyChatId: NewUUID(et.ReplyChatID),
				Text:        et.Text,
				Mentions:    NewUserMentions(et.Mentions),
			},
		}

	case *event.ChatEditEvent:
		return &EventStreamResponse_ChatEdit{
			ChatEdit: &ChatEditEvent{
				ChatId: NewUUID(et.ChatID),
				Text:   et.Text,
			},
		}

	default:
		panic("apiv1.NewEventStreamResponseEvent: " + reflect.TypeOf(et).String() + " is not implemented")
	}
	return nil
}
