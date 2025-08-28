// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatroom

import (
	"time"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv"
	"github.com/go-pogo/webapp"
	"github.com/go-pogo/webapp/logger"
	"github.com/roeldev/demo-chatroom/api/v1/apiv1connect"
	"github.com/roeldev/demo-chatroom/chatauth"
	"github.com/roeldev/demo-chatroom/chatevents"
	"github.com/roeldev/demo-chatroom/chatusers"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

type Config struct {
	Logger                 logger.Config       `env:",include"`
	Server                 webapp.ServerConfig `env:",include"`
	AllowedOrigins         []string            `env:"CORS_ALLOW_ORIGINS"`
	TypingIndicatorTimeout time.Duration       `default:"5s"`
}

var _ serv.RoutesRegisterer = (*Service)(nil)

type Service struct {
	log     zerolog.Logger
	conf    Config
	auth    chatauth.SignerParser
	manager *chatauth.Manager
	history *chatevents.HistoryHandler
	broker  *chatevents.EventsBroker
	users   chatusers.UsersStore

	interceptor connect.Interceptor
	cors        *cors.Cors
}

func NewService(conf Config, log *logger.Logger, opts ...Option) (*Service, error) {
	var err error
	svc := &Service{
		log:  log.Logger,
		conf: conf,
		cors: cors.New(cors.Options{
			AllowedOrigins: conf.AllowedOrigins,
			AllowedMethods: connectcors.AllowedMethods(),
			AllowedHeaders: append([]string{"authorization"}, connectcors.AllowedHeaders()...),
			ExposedHeaders: connectcors.ExposedHeaders(),
		}),
	}

	if err = svc.applyOpts(opts); err != nil {
		return nil, err
	}

	// todo: auth via priv/pub certs
	if svc.auth, err = chatauth.NewJWTAuth(); err != nil {
		return nil, err
	}
	if svc.users == nil {
		svc.users = chatusers.NewUsersStore(8)
	}
	if svc.history == nil {
		svc.history = chatevents.NewHistoryHandler(chatevents.NewEventsStore(32))
		if svc.broker != nil {
			svc.broker.Handle(svc.history)
		}
	}
	if svc.broker == nil {
		svc.broker = chatevents.NewEventsBroker(svc.log, svc.history)
	}

	svc.manager = chatauth.NewManager(svc.auth, svc.users, svc.broker)
	svc.interceptor = apiv1connect.NewHandlerInterceptor(svc.log, svc.auth, svc.users)
	return svc, nil
}

func (svc *Service) applyOpts(opts []Option) error {
	var err error
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		err = errors.Append(err, opt(svc))
	}
	return err
}

func (svc *Service) RegisterRoutes(rh serv.RouteHandler) {
	routes := []serv.Route{
		svc.authService(),
		svc.registryService(),
		svc.userService(),
		svc.eventsService(),
	}
	for _, route := range routes {
		route.Handler = svc.cors.Handler(route.Handler)
		rh.HandleRoute(route)
	}
}

func (svc *Service) authService() serv.Route {
	path, handler := apiv1connect.NewAuthServiceHandler(
		apiv1connect.NewAuthService(svc.log, svc.manager),
		connect.WithInterceptors(svc.interceptor),
	)
	return serv.Route{
		Name:    "auth-service",
		Pattern: path,
		Handler: handler,
	}
}

func (svc *Service) registryService() serv.Route {
	path, handler := apiv1connect.NewRegistryServiceHandler(
		apiv1connect.NewRegistryService(svc.log, svc.users),
		connect.WithInterceptors(svc.interceptor),
	)
	return serv.Route{
		Name:    "registry-service",
		Pattern: path,
		Handler: handler,
	}
}

func (svc *Service) userService() serv.Route {
	path, handler := apiv1connect.NewUserServiceHandler(
		apiv1connect.NewUserService(
			svc.log,
			svc.users,
			chatusers.NewTypingIndicator(svc.conf.TypingIndicatorTimeout),
			svc.broker,
		),
		connect.WithInterceptors(svc.interceptor),
	)
	return serv.Route{
		Name:    "user-service",
		Pattern: path,
		Handler: handler,
	}
}

func (svc *Service) eventsService() serv.Route {
	path, handler := apiv1connect.NewEventsServiceHandler(
		apiv1connect.NewEventsService(
			svc.log,
			svc.history,
			svc.broker,
			svc.manager,
		),
		connect.WithInterceptors(svc.interceptor),
	)
	return serv.Route{
		Name:    "events-service",
		Pattern: path,
		Handler: handler,
	}
}
