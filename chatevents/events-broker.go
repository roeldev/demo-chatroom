// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatevents

import (
	"sync"
	"time"

	"github.com/roeldev/demo-chatroom/chatevents/event"
	"github.com/rs/zerolog"
)

type EventHandler interface {
	HandleEvent(e Event)
}

type EventHandlerFunc func(e Event)

func (fn EventHandlerFunc) HandleEvent(e Event) { fn(e) }

type Publisher interface {
	Publish(typ event.Type)
}

var _ Publisher = (*EventsBroker)(nil)

type EventsBroker struct {
	log      zerolog.Logger
	mut      sync.RWMutex
	handlers []EventHandler
}

func NewEventsBroker(log zerolog.Logger, h ...EventHandler) *EventsBroker {
	if h == nil {
		h = make([]EventHandler, 0, 2)
	}

	return &EventsBroker{
		log:      log,
		handlers: h,
	}
}

func (eb *EventsBroker) Handle(h EventHandler) {
	eb.mut.Lock()
	defer eb.mut.Unlock()
	eb.handlers = append(eb.handlers, h)
}

func (eb *EventsBroker) Publish(typ event.Type) {
	e := Event{
		Time: time.Now(),
		Type: typ,
	}
	eb.log.Info().EmbedObject(e).Msg("publish")

	eb.mut.RLock()
	defer eb.mut.RUnlock()

	for _, h := range eb.handlers {
		go h.HandleEvent(e)
	}
}
