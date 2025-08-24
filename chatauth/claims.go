// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatauth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/roeldev/demo-chatroom/chatusers"
)

type Claims struct {
	jwt.RegisteredClaims

	UserID chatusers.UserID `json:"uid"`
}

func NewClaims(uid chatusers.UserID) *Claims {
	return &Claims{
		UserID: uid,
	}
}

func RenewClaims(old Claims) *Claims {
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    old.RegisteredClaims.Issuer,
			Subject:   old.RegisteredClaims.Subject,
			Audience:  old.RegisteredClaims.Audience,
			NotBefore: old.ExpiresAt,
			ID:        old.RegisteredClaims.ID,
		},

		UserID: old.UserID,
	}
}
