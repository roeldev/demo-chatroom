// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1connect

import (
	"context"
	"net"
	"strings"

	"connectrpc.com/connect"
	"github.com/go-pogo/errors"
	"github.com/roeldev/demo-chatroom/chatauth"
	"github.com/roeldev/demo-chatroom/chatusers"
	"github.com/rs/zerolog"
)

type handlerInterceptor struct {
	log    zerolog.Logger
	parser chatauth.Parser
	users  chatusers.UsersStore
}

func NewHandlerInterceptor(log zerolog.Logger, auth chatauth.Parser, users chatusers.UsersStore) connect.Interceptor {
	return &handlerInterceptor{
		log:    log,
		parser: auth,
		users:  users,
	}
}

func (hi *handlerInterceptor) authorize(ctx context.Context, bearer string, salter chatauth.SecretSalter) (context.Context, error) {
	if strings.HasPrefix(bearer, "bearer ") || strings.HasPrefix(bearer, "Bearer ") {
		bearer = bearer[7:]
	}

	claims, err := hi.parser.Parse(bearer, salter)
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidToken)
	}

	user, err := hi.users.Get(claims.UserID)
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidUserID)
	}

	ctx = context.WithValue(ctx, claimsKey{}, claims)
	ctx = context.WithValue(ctx, userKey{}, knownUser{
		ID:   claims.UserID,
		User: user,
	})
	return ctx, nil
}

func (hi *handlerInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if req.Spec().IsClient {
			panic("apiv1connect: Use NewClientInterceptor() instead")
		}

		var err error
		if bearer := req.Header().Get("authorization"); bearer != "" {
			ctx, err = hi.authorize(ctx, bearer, requestSecretSalter(req))
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
		}

		resp, err := next(ctx, req)
		if err != nil {
			hi.log.Warn().
				Err(err).
				Msg("error")
		}

		return resp, err
	}
}

func (hi *handlerInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		ctx, err := hi.authorize(ctx,
			conn.RequestHeader().Get("authorization"),
			streamSecretSalter(conn),
		)
		if err != nil {
			return connect.NewError(connect.CodeUnauthenticated, err)
		}

		err = next(ctx, conn)
		if err != nil {
			hi.log.Warn().
				Err(err).
				Msg("stream error")
		}
		return err
	}
}

func (hi *handlerInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

var _ chatauth.SecretSalter = (*secretSalter)(nil)

type secretSalter struct {
	IPAddr    string
	UserAgent string
}

func requestSecretSalter(req connect.AnyRequest) chatauth.SecretSalter {
	return secretSalter{
		IPAddr:    removePort(req.Peer().Addr),
		UserAgent: req.Header().Get("user-agent"),
	}
}

func streamSecretSalter(conn connect.StreamingHandlerConn) chatauth.SecretSalter {
	return secretSalter{
		IPAddr:    removePort(conn.Peer().Addr),
		UserAgent: conn.RequestHeader().Get("user-agent"),
	}
}

func removePort(addr string) string {
	host, _, err := net.SplitHostPort(addr)

	var addrErr *net.AddrError
	if errors.As(err, &addrErr) {
		return addrErr.Addr
	}
	return host
}

const separator byte = 33

func (s secretSalter) SaltSecret(secret []byte) []byte {
	secret = append(secret, separator)
	secret = append(secret, s.IPAddr...)
	secret = append(secret, separator)
	secret = append(secret, s.UserAgent...)

	return secret
}
