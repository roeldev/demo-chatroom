// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package welcomebot

import (
	"context"

	"connectrpc.com/connect"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/webapp/logger"
	apiv1 "github.com/roeldev/demo-chatroom/api/v1"
	"github.com/roeldev/demo-chatroom/api/v1/apiv1connect"
	"github.com/roeldev/demo-chatroom/chatbot"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Config struct {
	Logger logger.Config  `env:",include"`
	Bot    chatbot.Config `env:",include"`
}

type botClient = chatbot.Client

type WelcomeBot struct {
	*botClient
	events apiv1connect.EventsServiceClient

	log zerolog.Logger
}

func New(conf chatbot.Config, log zerolog.Logger) *WelcomeBot {
	bot := chatbot.New(conf)

	return &WelcomeBot{
		log:       log,
		botClient: bot,
		events: apiv1connect.NewEventsServiceClient(
			bot.HTTPClient(),
			conf.BaseURL(),
			connect.WithInterceptors(bot.Interceptor()),
		),
	}
}

func (bot *WelcomeBot) Login(ctx context.Context, user *apiv1.UserDetails) error {
	if err := bot.botClient.Login(ctx, user); err != nil {
		return errors.Wrap(err, "failed to join")
	}

	if _, err := bot.SendChat(ctx, connect.NewRequest(&apiv1.SendChatRequest{
		Time: timestamppb.Now(),
		Text: "Hi all!",
	})); err != nil {
		return errors.Wrap(err, "failed to send generic welcome chat")
	}
	return nil
}

func (bot *WelcomeBot) WelcomeUser(ctx context.Context, uid string, name string) error {
	bot.log.Debug().
		Str("user_id", uid).
		Str("user_name", name).
		Msg("send welcome chat")

	_, err := bot.SendChat(ctx, connect.NewRequest(&apiv1.SendChatRequest{
		Time: timestamppb.Now(),
		Text: "Welcome @" + name,
		Mentions: []*apiv1.UserMention{{
			UserId:   &apiv1.UUID{Value: uid},
			UserName: name,
		}},
	}))
	if err != nil {
		return errors.Wrap(err, "failed to welcome user")
	}
	return nil
}

func (bot *WelcomeBot) ListenForEvents(ctx context.Context) error {
	bot.log.Debug().Msg("start listening for events")
	stream, err := bot.events.EventStream(ctx, connect.NewRequest(&emptypb.Empty{}))
	if err != nil {
		return err
	}

	bot.log.Debug().Msg("listening for events...")
	for stream.Receive() {
		msg := stream.Msg()
		if join := msg.GetUserJoin(); join != nil {
			go func() {
				if err := bot.WelcomeUser(ctx,
					join.User.Id.Value,
					join.User.Details.Name,
				); err != nil {
					bot.log.Warn().Err(err).Msg("failed to welcome user")
				}
			}()
		}
	}
	return stream.Err()
}
