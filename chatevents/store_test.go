// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatevents

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLimitedEventsStore(t *testing.T) {
	assert.Equal(t,
		LimitedEventsStore{events: make([]Event, 0, defaultLimitedSize)},
		*NewLimitedEventsStore(0),
	)
}

func TestLimitedEventsStore_Cap(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var store LimitedEventsStore
		assert.Equal(t, 0, store.Cap())
	})
	t.Run("new", func(t *testing.T) {
		store := NewLimitedEventsStore(defaultLimitedSize)
		assert.Equal(t, defaultLimitedSize, store.Cap())
	})
}

func TestLimitedEventsStore_Len(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var store LimitedEventsStore
		assert.Equal(t, 0, store.Len())
	})
	t.Run("new", func(t *testing.T) {
		store := NewLimitedEventsStore(defaultLimitedSize)
		assert.Equal(t, 0, store.Len())
	})

}

func TestLimitedEventsStore_All(t *testing.T) {
	t.Run("add 1", func(t *testing.T) {
		store := NewLimitedEventsStore(defaultLimitedSize)
		first := Event{Time: time.Now()}
		store.Add(first)

		assert.Equal(t, 1, store.next)
		assert.Equal(t, []Event{first}, store.All())
	})
	t.Run("add max", func(t *testing.T) {
		store := NewLimitedEventsStore(defaultLimitedSize)
		for i := 0; i < defaultLimitedSize; i++ {
			store.Add(Event{})
		}
		assert.Equal(t, 0, store.next)
	})
	t.Run("add max + n", func(t *testing.T) {
		store := NewLimitedEventsStore(defaultLimitedSize)

		n := 1 + rand.IntN(4)
		for i := 0; i < n; i++ {
			store.Add(Event{})
		}

		want := make([]Event, defaultLimitedSize)
		for i := 0; i < defaultLimitedSize; i++ {
			want[i].Time = time.Date(2000+i, 1, 1, 1, 1, 1, 1, time.UTC)
			store.Add(want[i])
		}

		assert.Equal(t, n, store.next)
		assert.Equal(t, want, store.All())
	})
}

func TestLimitedEventsStore_ListEvents(t *testing.T) {
	store := NewLimitedEventsStore(defaultLimitedSize)

	n := 1 + rand.IntN(4)
	for i := 0; i < n; i++ {
		store.Add(Event{})
	}

	until := time.Date(2004, 11, 1, 1, 1, 1, 1, time.UTC)
	want := make([]Event, 0)
	for i := 0; i < defaultLimitedSize; i++ {
		e := Event{Time: time.Date(2000+i, 1, 1, 1, 1, 1, 1, time.UTC)}
		if e.Time.After(until) {
			want = append([]Event{e}, want...)
		}
		store.Add(e)
	}

	assert.Equal(t, n, store.next)
	assert.Equal(t, want, store.ListEvents(until, defaultLimitedSize))
}
