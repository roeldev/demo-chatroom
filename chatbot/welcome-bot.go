// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatbot

import (
	"context"

	"connectrpc.com/connect"
	"github.com/go-pogo/errors"
	chatroom "github.com/roeldev/demo-chatroom"
	apiv1 "github.com/roeldev/demo-chatroom/api/v1"
	"github.com/roeldev/demo-chatroom/api/v1/apiv1connect"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WelcomeBot struct {
	*botClient
	events apiv1connect.EventsServiceClient

	log zerolog.Logger
}

func NewWelcomeBot(conf chatroom.ClientConfig, log zerolog.Logger) *WelcomeBot {
	client := chatroom.NewClient(conf)
	return &WelcomeBot{
		log:       log,
		botClient: client,
		events: apiv1connect.NewEventsServiceClient(
			client.HTTPClient(),
			client.BaseURL(),
			connect.WithInterceptors(client.Interceptor()),
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

func (bot *WelcomeBot) Logout(ctx context.Context) (err error) {
	defer errors.AppendFunc(&err, func() error {
		return bot.botClient.Logout(ctx)
	})

	if _, err = bot.SendChat(ctx, connect.NewRequest(&apiv1.SendChatRequest{
		Time: timestamppb.Now(),
		Text: "Goodbye all!",
	})); err != nil {
		err = errors.Wrap(err, "failed to send generic goodbye chat")
	}
	return err
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
		return errors.Wrap(err, "failed to open stream")
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

func (bot *WelcomeBot) Run(ctx context.Context) error {
	ctx, cancelListen := context.WithCancelCause(ctx)

	go func() {
		defer cancelListen(nil)
		if err := bot.ListenForEvents(ctx); err != nil && !errors.Is(err, context.Canceled) {
			cancelListen(err)
		}
	}()

	<-ctx.Done()
	if err := bot.Logout(context.Background()); err != nil {
		return errors.Append(
			context.Cause(ctx),
			errors.Wrap(err, "failed to gracefully leave"),
		)
	}

	return context.Cause(ctx)
}
