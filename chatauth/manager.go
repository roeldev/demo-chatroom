// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatauth

import (
	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/roeldev/demo-chatroom/chatevents"
	"github.com/roeldev/demo-chatroom/chatevents/event"
	"github.com/roeldev/demo-chatroom/chatusers"
)

type Leaver interface {
	Leave(uid chatusers.UserID, reason event.LeaveReason)
}

type Manager struct {
	Signer
	users chatusers.UsersStore
	event chatevents.Publisher
}

func NewManager(signer Signer, users chatusers.UsersStore, pub chatevents.Publisher) *Manager {
	return &Manager{
		Signer: signer,
		users:  users,
		event:  pub,
	}
}

func (man *Manager) Join(user chatusers.User, salter SecretSalter) (chatusers.UserID, string, error) {
	uid, err := man.users.Add(user)
	if err != nil {
		return uuid.Nil, "", err
	}

	token, err := man.Signer.Sign(NewClaims(uid), salter)
	if err != nil {
		return uid, "", connect.NewError(connect.CodeInternal, err)
	}

	man.event.Publish(&event.UserJoinEvent{
		UserID:      uid,
		UserDetails: user.UserDetails,
		UserFlags:   user.Flags,
	})
	return uid, token, nil
}

func (man *Manager) Leave(uid chatusers.UserID, reason event.LeaveReason) {
	if user, ok := man.users.Delete(uid); ok {
		man.event.Publish(&event.UserLeaveEvent{
			UserID:      uid,
			UserDetails: user.UserDetails,
			Reason:      reason,
		})
	}
}
