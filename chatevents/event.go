// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatevents

import (
	"time"

	"github.com/roeldev/demo-chatroom/chatevents/event"
	"github.com/roeldev/demo-chatroom/chatusers"
	"github.com/rs/zerolog"
)

type EventType = event.Type

//type ChatID = uuid.UUID

var _ zerolog.LogObjectMarshaler = (*Event)(nil)

type Event struct {
	Time time.Time
	Type EventType
}

func (e Event) MarshalZerologObject(zl *zerolog.Event) {
	zl.Time("etime", e.Time)
	if obj, ok := e.Type.(zerolog.LogObjectMarshaler); ok {
		zl.EmbedObject(obj)
		return
	}

	zl.Type("type", e.Type)
	if ue, ok := e.Type.(event.UserEvent); ok {
		zl.Str("user", chatusers.IdentifierString(ue.GetUserID(), ue.GetUserDetails()))
	}
}
