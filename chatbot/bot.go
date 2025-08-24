// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatbot

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

type Config struct {
	ServerHost string             `env:"API_SERVER_HOST" default:"localhost"`
	ServerPort serv.Port          `env:"API_SERVER_PORT" default:"8080"`
	HTTPClient connect.HTTPClient `env:"-"`
}

func (c Config) BaseURL() string {
	// todo: tls support
	return "http://" + serv.JoinHostPort(c.ServerHost, c.ServerPort)
}

type Client struct {
	authServiceClient
	userServiceClient

	httpClient  connect.HTTPClient
	interceptor connect.Interceptor
}

func New(conf Config) *Client {
	if conf.HTTPClient == nil {
		conf.HTTPClient = http.DefaultClient
	}

	baseURL := conf.BaseURL()
	interceptor := apiv1connect.NewClientInterceptor("")
	interceptors := connect.WithInterceptors(interceptor)

	return &Client{
		httpClient:  conf.HTTPClient,
		interceptor: interceptor,

		authServiceClient: apiv1connect.NewAuthServiceClient(
			conf.HTTPClient,
			baseURL,
			interceptors,
		),
		userServiceClient: apiv1connect.NewUserServiceClient(
			conf.HTTPClient,
			baseURL,
			interceptors,
		),
	}
}

func (bot *Client) HTTPClient() connect.HTTPClient { return bot.httpClient }

func (bot *Client) Interceptor() connect.Interceptor { return bot.interceptor }

func (bot *Client) Login(ctx context.Context, user *apiv1.UserDetails) error {
	if _, err := bot.Join(ctx, connect.NewRequest(&apiv1.JoinRequest{
		User:  user,
		Flags: apiv1.UserFlag_IsBot,
	})); err != nil {
		return errors.Wrap(err, "failed to join")
	}
	return nil
}

func (bot *Client) Logout(ctx context.Context) error {
	if _, err := bot.Leave(ctx, connect.NewRequest(&emptypb.Empty{})); err != nil {
		return errors.Wrap(err, "failed to gracefully leave")
	}
	return nil
}
