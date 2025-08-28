// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatevents

import (
	"time"

	"github.com/google/uuid"
	"github.com/roeldev/demo-chatroom/chatevents/event"
	"github.com/roeldev/demo-chatroom/chatusers"
	"github.com/rs/zerolog"
)

var _ zerolog.LogObjectMarshaler = (*Event)(nil)

type Event struct {
	Time time.Time
	Type event.Type
}

// AsUserEvent casts Type as a [event.UserEvent] and returns it or nil.
func (e Event) AsUserEvent() event.UserEvent {
	if ue, ok := e.Type.(event.UserEvent); ok {
		return ue
	}
	return nil
}

// AsReceiverEvent casts Type as a [event.ReceiverEvent] and returns it or nil.
func (e Event) AsReceiverEvent() event.ReceiverEvent {
	if re, ok := e.Type.(event.ReceiverEvent); ok {
		return re
	}
	return nil
}

func (e Event) MarshalZerologObject(ze *zerolog.Event) {
	ze.Time("etime", e.Time)
	if obj, ok := e.Type.(zerolog.LogObjectMarshaler); ok {
		ze.EmbedObject(obj)
		return
	}

	ze.Type("type", e.Type)
	if ue := e.AsUserEvent(); ue != nil {
		ze.Str("user", chatusers.IdentifierString(ue.GetUserID(), ue.GetUserDetails()))
	}
	if re := e.AsReceiverEvent(); re != nil {
		if rid := re.GetReceiverID(); rid != uuid.Nil {
			ze.Stringer("receiver_id", rid)
		}
	}
}
