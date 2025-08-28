// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatevents

import (
	"sync"
	"time"

	"github.com/roeldev/demo-chatroom/chatevents/event"
)

type EventsLister interface {
	ListEvents(until time.Time, limit uint64) []Event
}

type EventsStore interface {
	EventsLister
	All() []Event
	Add(e Event)
	UpdateChatEvent(id event.ChatID, fn func(*event.ChatEvent))
}

type eventsStore struct {
	mut         sync.RWMutex
	events      []Event
	index, size uint8
}

// NewEventsStore returns a simple in-memory [EventsStore] with a specified
// size limit.
func NewEventsStore(size uint8) EventsStore {
	return &eventsStore{
		events: make([]Event, size),
	}
}

func (es *eventsStore) All() []Event {
	es.mut.RLock()
	defer es.mut.RUnlock()

	clone := make([]Event, len(es.events))
	for i, v := range es.events {
		j := uint8(i)
		if j < es.index {
			j += es.index
		} else {
			j -= es.index
		}
		clone[j] = v
	}
	return clone
}

func (es *eventsStore) ListEvents(until time.Time, limit uint64) []Event {
	//TODO implement me
	panic("implement me")
}

func (es *eventsStore) Add(e Event) {
	es.mut.Lock()
	defer es.mut.Unlock()

	es.events[es.index] = e
	es.index++
	if es.index >= es.size {
		es.index = 0
	}
}

func (es *eventsStore) UpdateChatEvent(id event.ChatID, fn func(*event.ChatEvent)) {
	es.mut.Lock()
	defer es.mut.Unlock()

	for _, e := range es.events {
		if chat, ok := e.Type.(*event.ChatEvent); ok && chat.ChatID == id {
			fn(chat)
			return
		}
	}
}
