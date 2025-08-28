// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1connect

import (
	"context"

	"connectrpc.com/connect"
	"github.com/go-pogo/errors"
	"github.com/google/uuid"
	"github.com/roeldev/demo-chatroom/api/v1"
	"github.com/roeldev/demo-chatroom/chatevents"
	"github.com/roeldev/demo-chatroom/chatevents/event"
	"github.com/roeldev/demo-chatroom/chatusers"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ UserServiceHandler = (*UserService)(nil)

type UserService struct {
	log   zerolog.Logger
	users chatusers.UsersStore
	event chatevents.Publisher
}

func NewUserService(log zerolog.Logger, users chatusers.UsersStore, pub chatevents.Publisher) *UserService {
	return &UserService{
		log:   log,
		users: users,
		event: pub,
	}
}

func (svc *UserService) UpdateDetails(_ context.Context, req *connect.Request[apiv1.UpdateDetailsRequest]) (*connect.Response[emptypb.Empty], error) {
	//TODO implement me
	panic("implement me")
}

func (svc *UserService) UpdateStatus(ctx context.Context, req *connect.Request[apiv1.UpdateStatusRequest]) (*connect.Response[emptypb.Empty], error) {
	user := getUser(ctx)
	before := user.Status

	user.Status = req.Msg.Status.ToChatUserStatus()
	if err := svc.users.Update(user.ID, user.User); err != nil {
		return nil, errors.Wrap(err, ErrChangeUserStatus)
	}

	go svc.event.Publish(&event.UserStatusEvent{
		UserID:      user.ID,
		UserDetails: user.UserDetails,
		Before:      before,
		After:       user.Status,
	})

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (svc *UserService) IndicateTyping(ctx context.Context, req *connect.Request[apiv1.IndicateTypingRequest]) (*connect.Response[emptypb.Empty], error) {
	user := getUser(ctx)

	receiver, err := req.Msg.ReceiverId.ParseUUID()
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// todo: timer laten lopen die na ~5 sec. typing op false zet en dit published

	go svc.event.Publish(&event.UserTypingEvent{
		UserID:      user.ID,
		UserDetails: user.UserDetails,
		ReceiverID:  receiver,
		IsTyping:    req.Msg.Typing,
	})

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (svc *UserService) SendChat(ctx context.Context, req *connect.Request[apiv1.SendChatRequest]) (*connect.Response[emptypb.Empty], error) {
	user := getUser(ctx)

	receiver, err := req.Msg.ReceiverId.ParseUUID()
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	go svc.event.Publish(&event.ChatEvent{
		ChatID:      uuid.New(),
		UserID:      user.ID,
		UserDetails: user.UserDetails,
		ReceiverID:  receiver,
		Text:        req.Msg.Text,
	})

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (svc *UserService) EditChat(ctx context.Context, c *connect.Request[apiv1.ChatEditEvent]) (*connect.Response[emptypb.Empty], error) {
	//TODO implement me
	panic("implement me")
}

func (svc *UserService) EmojiReply(ctx context.Context, c *connect.Request[apiv1.EmojiReplyRequest]) (*connect.Response[emptypb.Empty], error) {
	//TODO implement me
	panic("implement me")
}
