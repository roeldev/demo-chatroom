// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1

import (
	"reflect"

	"github.com/roeldev/demo-chatroom/chatevents/event"
)

type PreviousEventsResponseEvent = isPreviousEventsResponse_PreviousEvent_Event

func NewPreviousEventsResponseEvent(typ event.Type) PreviousEventsResponseEvent {
	switch et := typ.(type) {
	default:
		panic("apiv1.NewPreviousEventsResponseEvent: " + reflect.TypeOf(et).String() + " is not implemented")
	}
	return nil
}
