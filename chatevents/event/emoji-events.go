// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package event

import "github.com/roeldev/demo-chatroom/chatusers"

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

func (e EmojiReplyEvent) GetUserID() chatusers.UserID           { return e.UserID }
func (e EmojiReplyEvent) GetUserDetails() chatusers.UserDetails { return e.UserDetails }
