// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatroom

import (
	"github.com/roeldev/demo-chatroom/chatusers"
)

type Option func(svc *Service) error

func WithUsers(users chatusers.UsersStore) Option {
	return func(svc *Service) error {
		svc.users = users
		return nil
	}
}
