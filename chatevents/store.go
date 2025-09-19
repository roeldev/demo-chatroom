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
	ListEvents(until time.Time, limit int) []Event
}

type EventsStore interface {
	EventsLister
	All() []Event
	Add(e Event)
	UpdateChatEvent(id event.ChatID, fn func(*event.ChatEvent))
}

const defaultLimitedSize = 32

var DefaultLimitedSize uint8 = defaultLimitedSize

type LimitedEventsStore struct {
	mut    sync.RWMutex
	events []Event
	next   int
}

// NewLimitedEventsStore returns a simple in-memory [EventsStore] with a
// specified size limit which is at least 8.
func NewLimitedEventsStore(size uint8) *LimitedEventsStore {
	var es LimitedEventsStore
	es.init(size)
	return &es
}

func (es *LimitedEventsStore) init(size uint8) {
	if size < 0 {
		size = defaultLimitedSize
	}

	es.events = make([]Event, 0, size)
}

func (es *LimitedEventsStore) Cap() int {
	es.mut.RLock()
	defer es.mut.RUnlock()
	return cap(es.events)
}

func (es *LimitedEventsStore) Len() int {
	es.mut.RLock()
	defer es.mut.RUnlock()
	return len(es.events)
}

func (es *LimitedEventsStore) All() []Event {
	es.mut.RLock()
	defer es.mut.RUnlock()

	size := len(es.events)
	clone := make([]Event, 0, size)

	for i := 0; i < size; i++ {
		j := i + es.next
		if j >= size {
			j -= size
		}
		clone = append(clone, es.events[j])
	}
	return clone
}

func (es *LimitedEventsStore) ListEvents(until time.Time, limit int) []Event {
	es.mut.RLock()
	defer es.mut.RUnlock()

	size := len(es.events)
	if limit > size || limit <= 0 {
		limit = size
	}

	clone := make([]Event, 0)
	for i := 0; i < limit; i++ {
		j := i + es.next
		if j >= size {
			j -= size
		}

		e := es.events[j]
		if e.Time.Before(until) {
			continue
		}

		clone = append([]Event{es.events[j]}, clone...)
	}
	return clone
}

func (es *LimitedEventsStore) Add(e Event) {
	es.mut.Lock()
	defer es.mut.Unlock()

	if es.events == nil {
		es.init(DefaultLimitedSize)
	}

	if len(es.events) < cap(es.events) {
		es.events = append(es.events, e)
	} else {
		es.events[es.next] = e
	}

	es.next++
	if es.next >= cap(es.events) {
		es.next = 0
	}
}

func (es *LimitedEventsStore) UpdateChatEvent(id event.ChatID, fn func(*event.ChatEvent)) {
	es.mut.Lock()
	defer es.mut.Unlock()

	for _, e := range es.events {
		if chat, ok := e.Type.(*event.ChatEvent); ok && chat.ChatID == id {
			fn(chat)
			return
		}
	}
}
