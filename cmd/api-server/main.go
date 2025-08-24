// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-pogo/env"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv"
	"github.com/go-pogo/webapp"
	"github.com/go-pogo/webapp/autoenv"
	"github.com/go-pogo/webapp/logger"
	"github.com/roeldev/demo-chatroom"
	"github.com/roeldev/demo-chatroom/chatauth"
	logpkg "github.com/rs/zerolog/log"
)

func main() {
	var conf chatroom.Config
	if err := unmarshal(&conf); err != nil {
		logpkg.Fatal().Err(err).Msg("failed to unmarshal config")
	}

	log := logger.New(conf.Logger)
	app, err := setup(conf, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup")
	}

	if err = webapp.Run(context.Background(),
		app.Run,
	); err != nil && !errors.Is(err, context.Canceled) {
		log.Err(err).Msg("error during execution")
	}

	// shutdown services
	if err = webapp.ShutdownTimeout(context.Background(), 10*time.Second,
		app.Shutdown,
	); err != nil {
		log.Err(err).Msg("error during shutdown")
	}

	log.Debug().Msg("goodbye")
}

func unmarshal(conf *chatroom.Config) error {
	loader := autoenv.NewLoader()
	if err := loader.Load(); err != nil {
		return err
	}

	if err := env.NewDecoder(env.System()).Decode(conf); err != nil {
		return err
	}

	if devBuild {
		_ = conf.Server.TLS.CACertFile.Set(loader.PrefixDir(conf.Server.TLS.CACertFile.String()))
		conf.Server.TLS.CertFile = loader.PrefixDir(conf.Server.TLS.CertFile)
		conf.Server.TLS.KeyFile = loader.PrefixDir(conf.Server.TLS.KeyFile)
	}
	return nil
}

func setup(conf chatroom.Config, log *logger.Logger) (*webapp.Base, error) {
	base, err := webapp.New("api-server",
		webapp.WithLogger(log),
		webapp.WithServerConfig(conf.Server),
		//webapp.WithServerOption(serv.WithAllowUnencryptedHTTP2(true)),
		//webapp.WithHealthChecker(),
		webapp.WithIgnoreFaviconRoute(),
	)
	if err != nil {
		return nil, errors.Wrap(errors.WithExitCode(err, 1), "unable to create webapp")
	}

	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// todo: wel toevoegen dmn http write context struct ding
	// https://pkg.go.dev/net/http#ResponseController

	server := base.Server()
	if serv.ShouldUseTLS(server.TLSConfig) {
		server.Protocols.SetUnencryptedHTTP2(true)
	}

	//base.Server().Config.ReadTimeout = 0
	server.Config.WriteTimeout = 0
	server.Config.IdleTimeout = 0
	server.Config.MaxHeaderBytes = http.DefaultMaxHeaderBytes

	// todo: auth via priv/pub certs
	auth, err := chatauth.NewJWTAuth()
	if err != nil {
		return nil, err
	}

	api, err := chatroom.NewService(conf, log, auth)
	if err != nil {
		return nil, err
	}

	api.RegisterRoutes(base.RouteHandler())
	return base, nil
}
