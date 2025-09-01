// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1connect

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/roeldev/demo-chatroom/api/v1"
	"github.com/roeldev/demo-chatroom/chatauth"
	"github.com/roeldev/demo-chatroom/chatevents/event"
	"github.com/roeldev/demo-chatroom/chatusers"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ AuthServiceHandler = (*AuthService)(nil)

type AuthService struct {
	log  zerolog.Logger
	auth *chatauth.Manager
}

func NewAuthService(log zerolog.Logger, auth *chatauth.Manager) *AuthService {
	return &AuthService{
		log:  log,
		auth: auth,
	}
}

func (svc *AuthService) Join(_ context.Context, req *connect.Request[apiv1.JoinRequest]) (*connect.Response[apiv1.JoinResponse], error) {
	user, err := chatusers.NewUser("", apiv1.FromJoinRequest(req.Msg))
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	uid, token, err := svc.auth.Join(*user, requestSecretSalter(req))
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	svc.log.Info().
		Str("user", chatusers.IdentifierString(uid, user.UserDetails)).
		Str("token", token).
		Msg("user joins")

	return connect.NewResponse(&apiv1.JoinResponse{
		Token: token,
	}), nil
}

func (svc *AuthService) Keepalive(ctx context.Context, stream *connect.ClientStream[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	streamStart := time.Now()

	user := getUser(ctx)
	svc.log.Debug().
		Stringer("user", user).
		Msg("start keepalive")

	defer func() {
		svc.auth.Leave(user.ID, event.Disconnected)

		svc.log.Debug().
			Stringer("user", user).
			Stringer("duration", time.Since(streamStart)).
			Msg("close keepalive")
	}()

	for stream.Receive() {
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (svc *AuthService) Renew(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[apiv1.RenewResponse], error) {
	oldClaims := getClaims(ctx)
	newClaims := chatauth.RenewClaims(oldClaims)

	token, err := svc.auth.Sign(newClaims, requestSecretSalter(req))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	user := getUser(ctx)
	svc.log.Info().
		Stringer("user", user).
		Str("token", token).
		Msg("renew token")

	return connect.NewResponse(&apiv1.RenewResponse{
		Token: token,
	}), nil
}

func (svc *AuthService) Leave(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	user := getUser(ctx)
	svc.log.Info().
		Stringer("user", user).
		Msg("user leaves")

	svc.auth.Leave(user.ID, event.UserLeave)
	return connect.NewResponse(&emptypb.Empty{}), nil
}
