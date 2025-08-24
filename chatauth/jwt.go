// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatauth

import (
	"time"

	"github.com/go-pogo/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var _ SignerParser = (*JWTAuth)(nil)

type JWTAuth struct {
	secret  []byte
	expires time.Duration
}

func NewJWTAuth() (*JWTAuth, error) {
	return &JWTAuth{
		expires: time.Hour,
	}, nil
}

func (auth *JWTAuth) Sign(claims *Claims, salter SecretSalter) (string, error) {
	if claims == nil {
		panic("chatauth.JWTAuth: claims must not be nil")
	}
	if claims.UserID == uuid.Nil {
		return "", errors.New("missing UserID")
	}

	if claims.RegisteredClaims.IssuedAt == nil {
		claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Now())
	}
	if auth.expires != 0 {
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(
			claims.RegisteredClaims.IssuedAt.Add(auth.expires),
		)
	}

	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(auth.saltSecret(salter))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return tok, nil
}

func (auth *JWTAuth) Parse(token string, salter SecretSalter) (Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(token, &claims,
		func(_ *jwt.Token) (any, error) { return auth.saltSecret(salter), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return claims, errors.WithStack(err)
	}

	if auth.expires > 0 && claims.ExpiresAt.Before(time.Now()) {
		return claims, errors.New("token time has expired")
	}

	return claims, nil
}

func (auth *JWTAuth) saltSecret(salter SecretSalter) []byte {
	if salter == nil {
		return auth.secret
	}
	return salter.SaltSecret(auth.secret)
}
