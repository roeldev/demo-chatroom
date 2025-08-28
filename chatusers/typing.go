// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatusers

import (
	"sync"
	"time"
)

type TypingIndicator interface {
	IndicateTyping(uid UserID, typing bool, callback func(typing bool))
}

type typingIndicator struct {
	mut     sync.RWMutex
	timeout time.Duration
	timers  map[UserID]*time.Timer
}

func NewTypingIndicator(timeout time.Duration) TypingIndicator {
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	return &typingIndicator{
		timeout: timeout,
		timers:  make(map[UserID]*time.Timer, 8),
	}
}

func (t *typingIndicator) timer(uid UserID) *time.Timer {
	t.mut.Lock()
	defer t.mut.Unlock()

	if _, has := t.timers[uid]; !has {
		t.timers[uid] = time.NewTimer(0)
	}
	return t.timers[uid]
}

const panicNilCallback = "chatusers: callback must not be nil!"

func (t *typingIndicator) IndicateTyping(uid UserID, typing bool, callback func(typing bool)) {
	if callback == nil {
		panic(panicNilCallback)
	}

	timer := t.timer(uid)
	timer.Stop()

	if !typing {
		go callback(false)
		return
	}

	go callback(true)
	timer.Reset(t.timeout)

	go func() {
		<-timer.C
		timer.Stop()
		callback(false)
	}()
}
