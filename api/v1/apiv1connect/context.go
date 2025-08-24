// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1connect

import (
	"context"

	"github.com/roeldev/demo-chatroom/chatauth"
	"github.com/roeldev/demo-chatroom/chatusers"
)

type claimsKey struct{}

func getClaims(ctx context.Context) chatauth.Claims {
	if v := ctx.Value(claimsKey{}); v != nil {
		return v.(chatauth.Claims)
	}
	return chatauth.Claims{}
}

type userKey struct{}

func getUser(ctx context.Context) knownUser {
	if v := ctx.Value(userKey{}); v != nil {
		return v.(knownUser)
	}
	return knownUser{}
}

type knownUser struct {
	chatusers.User
	ID chatusers.UserID
}

func (u knownUser) String() string {
	return chatusers.IdentifierString(u.ID, u.UserDetails)
}
