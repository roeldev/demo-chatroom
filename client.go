// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatroom

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv"
	apiv1 "github.com/roeldev/demo-chatroom/api/v1"
	"github.com/roeldev/demo-chatroom/api/v1/apiv1connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authServiceClient = apiv1connect.AuthServiceClient
type userServiceClient = apiv1connect.UserServiceClient

type ClientConfig struct {
	ServerHost string             `env:"API_SERVER_HOST" default:"localhost"`
	ServerPort serv.Port          `env:"API_SERVER_PORT" default:"8080"`
	HTTPClient connect.HTTPClient `env:"-"`
}

func (c ClientConfig) httpClient() connect.HTTPClient {
	if c.HTTPClient == nil {
		return http.DefaultClient
	}
	return c.HTTPClient
}

func (c ClientConfig) BaseURL() string {
	// todo: tls support
	return "http://" + serv.JoinHostPort(c.ServerHost, c.ServerPort)
}

type Client struct {
	authServiceClient
	userServiceClient

	httpClient  connect.HTTPClient
	interceptor connect.Interceptor
}

func NewClient(conf ClientConfig) *Client {
	client := conf.httpClient()
	baseURL := conf.BaseURL()
	interceptor := apiv1connect.NewClientInterceptor("")

	return &Client{
		httpClient:  client,
		interceptor: interceptor,

		authServiceClient: apiv1connect.NewAuthServiceClient(
			client,
			baseURL,
			//connect.WithGRPC(),
			connect.WithInterceptors(interceptor),
		),
		userServiceClient: apiv1connect.NewUserServiceClient(
			client,
			baseURL,
			//connect.WithGRPC(),
			connect.WithInterceptors(interceptor),
		),
	}
}

func (c *Client) HTTPClient() connect.HTTPClient { return c.httpClient }

func (c *Client) Interceptor() connect.Interceptor { return c.interceptor }

func (c *Client) Login(ctx context.Context, user *apiv1.UserDetails) error {
	if _, err := c.Join(ctx, connect.NewRequest(&apiv1.JoinRequest{
		User:  user,
		Flags: apiv1.UserFlag_USER_FLAG_IS_BOT,
	})); err != nil {
		return errors.Wrap(err, "failed to join")
	}
	return nil
}

func (c *Client) Logout(ctx context.Context) error {
	if _, err := c.Leave(ctx, connect.NewRequest(&emptypb.Empty{})); err != nil {
		return errors.Wrap(err, "failed to gracefully leave")
	}
	return nil
}
