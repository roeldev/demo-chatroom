// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1connect

import (
	"context"

	"connectrpc.com/connect"
	"github.com/roeldev/demo-chatroom/api/v1"
	"github.com/roeldev/demo-chatroom/chatauth"
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

	_, token, err := svc.auth.Join(*user, requestSecretSalter(req))
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&apiv1.JoinResponse{
		Token: token,
	}), nil
}

func (svc *AuthService) Keepalive(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	//TODO implement me
	panic("implement me")
}

func (svc *AuthService) Renew(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[apiv1.RenewResponse], error) {
	oldClaims := getClaims(ctx)
	newClaims := chatauth.RenewClaims(oldClaims)

	token, err := svc.auth.Sign(newClaims, requestSecretSalter(req))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&apiv1.RenewResponse{
		Token: token,
	}), nil
}

func (svc *AuthService) Leave(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	svc.auth.Leave(getClaims(ctx).UserID)
	return connect.NewResponse(&emptypb.Empty{}), nil
}
