// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/webapp"
	"github.com/go-pogo/webapp/autoenv"
	"github.com/go-pogo/webapp/logger"
	"github.com/roeldev/demo-chatroom/api/v1"
	"github.com/roeldev/demo-chatroom/chatbot"
	logpkg "github.com/rs/zerolog/log"
)

func main() {
	var conf chatbot.Config
	if err := autoenv.Unmarshal(&conf); err != nil {
		logpkg.Fatal().Err(err).Msg("failed to unmarshal config")
	}

	log := logger.New(conf.Logger)
	bot := chatbot.NewWelcomeBot(conf.Client, log.Logger)

	ctx := context.Background()
	if err := bot.Login(ctx, &apiv1.UserDetails{
		Name:     "Welcome-bot",
		Initials: "hi",
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to join")
	}

	log.Debug().Msg("welcome-bot joined the chatroom")
	if err := webapp.Run(ctx, bot.Run); err != nil && !errors.Is(err, context.Canceled) {
		log.Warn().Err(err).Msg("error during run")
	}

	log.Debug().Msg("goodbye")
}
