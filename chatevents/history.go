// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatevents

import (
	"github.com/google/uuid"
	"github.com/roeldev/demo-chatroom/chatevents/event"
)

var (
	_ EventsLister = (*HistoryHandler)(nil)
	_ EventHandler = (*HistoryHandler)(nil)
)

type HistoryHandler struct {
	EventsStore
}

func NewHistoryHandler(store EventsStore) *HistoryHandler {
	if store == nil {
		store = NewEventsStore(32)
	}
	return &HistoryHandler{
		EventsStore: store,
	}
}

func (his *HistoryHandler) HandleEvent(e Event) {
	if re := e.AsReceiverEvent(); re != nil && re.GetReceiverID() != uuid.Nil {
		// do not store events in history when sent to a receiver
		return
	}

	switch et := e.Type.(type) {
	case *event.UserTypingEvent:
		// user typing is not stored in history
		return

	case *event.ChatEditEvent:
		his.UpdateChatEvent(et.ChatID, func(chat *event.ChatEvent) {
			chat.SetEdited(e.Time, et.Text)
		})

	case *event.EmojiReplyEvent:
		his.UpdateChatEvent(et.ReplyChatID, func(chat *event.ChatEvent) {
			chat.AddEmojiReply(event.EmojiReply{
				Time:        e.Time,
				UserID:      et.UserID,
				UserDetails: et.UserDetails,
				Emoji:       et.Emoji,
			})
		})

	case *event.EmojiRemoveEvent:
		his.UpdateChatEvent(et.ReplyChatID, func(chat *event.ChatEvent) {
			chat.RemoveEmojiReply(et.UserID)
		})

	default:
		his.Add(e)
	}
}
