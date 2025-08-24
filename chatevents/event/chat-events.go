// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/roeldev/demo-chatroom/chatusers"
	"github.com/rs/zerolog"
)

type ChatID = uuid.UUID

type ChatEdit struct {
	Time     time.Time
	Original string
}

type EmojiReply struct {
	Time        time.Time
	UserID      chatusers.UserID
	UserDetails chatusers.UserDetails
	Emoji       string
}

var _ zerolog.LogObjectMarshaler = (*ChatEvent)(nil)

type ChatEvent struct {
	event
	ChatID       ChatID
	UserID       chatusers.UserID
	UserDetails  chatusers.UserDetails
	ReceiverID   chatusers.UserID
	ReplyChatID  ChatID
	Text         string
	Edit         *ChatEdit
	Mentions     map[chatusers.UserID]string
	EmojiReplies map[chatusers.UserID]EmojiReply
}

func (e *ChatEvent) MarshalZerologObject(ze *zerolog.Event) {
	ze.Stringer("chat_id", e.ChatID)
	ze.Str("sender", chatusers.IdentifierString(e.UserID, e.UserDetails))
	if e.ReceiverID != uuid.Nil {
		ze.Stringer("receiver_id", e.ReceiverID)
	}
	if e.ReplyChatID != uuid.Nil {
		ze.Stringer("reply_chat_id", e.ReplyChatID)
	}
	ze.Str("text", e.Text)
}

func (e *ChatEvent) SetEdited(editTime time.Time, text string) {
	if e.Edit == nil {
		e.Edit = &ChatEdit{Original: e.Text}
	}

	e.Edit.Time = editTime
	e.Text = text
}

func (e *ChatEvent) AddEmojiReply(er EmojiReply) {
	if e.EmojiReplies == nil {
		e.EmojiReplies = make(map[chatusers.UserID]EmojiReply)
	}
	e.EmojiReplies[er.UserID] = er
}

func (e *ChatEvent) RemoveEmojiReply(uid chatusers.UserID) {
	if e.EmojiReplies != nil {
		return
	}
	delete(e.EmojiReplies, uid)
}

type ChatEditEvent struct {
	event
	ChatID ChatID
	Text   string
}

type EmojiReplyEvent struct {
	event
	UserID      chatusers.UserID
	UserDetails chatusers.UserDetails
	ReplyChatID ChatID
	Emoji       string
}

type EmojiRemoveEvent struct {
	event
	UserID      chatusers.UserID
	ReplyChatID ChatID
}

func (e *ChatEvent) GetUserID() chatusers.UserID      { return e.UserID }
func (e EmojiReplyEvent) GetUserID() chatusers.UserID { return e.UserID }

func (e *ChatEvent) GetUserDetails() chatusers.UserDetails      { return e.UserDetails }
func (e EmojiReplyEvent) GetUserDetails() chatusers.UserDetails { return e.UserDetails }
