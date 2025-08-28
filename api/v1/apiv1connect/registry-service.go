// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1connect

import (
	"context"

	"connectrpc.com/connect"
	"github.com/roeldev/demo-chatroom/api/v1"
	"github.com/roeldev/demo-chatroom/chatusers"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ RegistryServiceHandler = (*RegistryService)(nil)

type RegistryService struct {
	log   zerolog.Logger
	users chatusers.UsersStore
}

func NewRegistryService(log zerolog.Logger, users chatusers.UsersStore) *RegistryService {
	return &RegistryService{
		log:   log,
		users: users,
	}
}

func (svc *RegistryService) ActiveUsers(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[apiv1.ActiveUsersResponse], error) {
	users := svc.users.All()
	response := make([]*apiv1.ActiveUsersResponse_User, 0, len(users))

	for uid, user := range users {
		response = append(response, &apiv1.ActiveUsersResponse_User{
			Id:      apiv1.NewUUID(uid),
			Details: apiv1.NewUserDetails(user.UserDetails),
			Flags:   apiv1.NewUserFlags(user.Flags),
			Status:  apiv1.NewUserStatus(user.Status),
		})
	}

	return connect.NewResponse(&apiv1.ActiveUsersResponse{
		Time:  timestamppb.Now(),
		Users: response,
	}), nil
}
