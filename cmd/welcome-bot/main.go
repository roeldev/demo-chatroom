// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/webapp"
	"github.com/go-pogo/webapp/logger"
	"github.com/roeldev/demo-chatroom/api/v1"
	"github.com/roeldev/demo-chatroom/app/welcome-bot"
)

func main() {
	var conf welcomebot.Config
	conf.Bot.ServerHost = "localhost"
	conf.Bot.ServerPort = 8080

	log := logger.New(conf.Logger)
	bot := welcomebot.New(conf.Bot, log.Logger)

	ctx := context.Background()
	if err := bot.Login(ctx, &apiv1.UserDetails{
		Name:     "Welcome-bot",
		Initials: "hi",
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to join")
	}

	log.Debug().Msg("welcome-bot joined the chatroom")

	if err := webapp.Run(ctx, bot.ListenForEvents); err != nil && !errors.Is(err, context.Canceled) {
		log.Warn().Err(err).Msg("error while listening for events")
	}

	if err := bot.Logout(context.Background()); err != nil {
		log.Warn().Err(err).Msg("failed to gracefully leave")
	}
	log.Debug().Msg("goodbye")
}
