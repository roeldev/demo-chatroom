// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1

import (
	"reflect"

	"github.com/roeldev/demo-chatroom/chatevents/event"
)

func NewEventUser(u event.UserEvent) *EventUser {
	return &EventUser{
		Id:      NewUUID(u.GetUserID()),
		Details: NewUserDetails(u.GetUserDetails()),
	}
}

type EventStreamResponseEvent = isEventStreamResponse_Event

// NewEventStreamResponseEvent translates an [event.Type] to an
// [EventStreamResponseEvent].
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
				User:   NewEventUser(et),
				Reason: NewLeaveReason(et.Reason),
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
			ChatSent: NewChatSentEvent(et),
		}

	//case *event.ChatEditEvent:
	//	return &EventStreamResponse_ChatEdit{
	//		ChatEdit: &ChatEditEvent{
	//			ChatId: NewUUID(et.ChatID),
	//			Text:   et.Text,
	//		},
	//	}

	//case *event.EmojiReplyEvent:

	default:
		panic("apiv1.NewEventStreamResponseEvent: " + reflect.TypeOf(et).String() + " is not implemented")
	}
	return nil
}

type PreviousEventsResponseEvent = isPreviousEventsResponse_PreviousEvent_Event

// NewPreviousEventsResponseEvent translates an [event.Type] to a
// [PreviousEventsResponseEvent].
func NewPreviousEventsResponseEvent(typ event.Type) PreviousEventsResponseEvent {
	switch et := typ.(type) {
	case *event.UserJoinEvent:
		return &PreviousEventsResponse_PreviousEvent_UserJoin{
			UserJoin: &UserJoinEvent{
				User:  NewEventUser(et),
				Flags: NewUserFlags(et.UserFlags),
			},
		}

	case *event.UserLeaveEvent:
		return &PreviousEventsResponse_PreviousEvent_UserLeave{
			UserLeave: &UserLeaveEvent{
				User:   NewEventUser(et),
				Reason: NewLeaveReason(et.Reason),
			},
		}

	case *event.UserUpdateEvent:
		return &PreviousEventsResponse_PreviousEvent_UserUpdate{
			UserUpdate: &UserUpdateEvent{
				User:   NewEventUser(et),
				Before: NewUserDetails(et.Before),
			},
		}

	case *event.ChatEvent:
		return &PreviousEventsResponse_PreviousEvent_ChatSent{
			ChatSent: NewChatSentEvent(et),
		}

	default:
		panic("apiv1.NewPreviousEventsResponseEvent: " + reflect.TypeOf(et).String() + " is not implemented")
	}
	return nil
}
