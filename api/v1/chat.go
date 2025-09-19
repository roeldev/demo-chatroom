// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1

import (
	"github.com/roeldev/demo-chatroom/chatevents/event"
	"github.com/roeldev/demo-chatroom/chatusers"
)

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

func NewChatSentEvent(et *event.ChatEvent) *ChatSentEvent {
	return &ChatSentEvent{
		ChatId:      NewUUID(et.ChatID),
		User:        NewEventUser(et),
		ReceiverId:  NewUUID(et.ReceiverID),
		ReplyChatId: NewUUID(et.ReplyChatID),
		Text:        et.Text,
		Mentions:    NewUserMentions(et.Mentions),
	}
}
