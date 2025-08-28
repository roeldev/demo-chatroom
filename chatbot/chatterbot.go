// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatbot

import (
	"context"
	"math/rand/v2"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/go-faker/faker/v4"
	chatroom "github.com/roeldev/demo-chatroom"
	apiv1 "github.com/roeldev/demo-chatroom/api/v1"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatterBot struct {
	*botClient

	log  zerolog.Logger
	rand random
}

func NewChatterBot(conf chatroom.ClientConfig, log zerolog.Logger) *ChatterBot {
	return &ChatterBot{
		botClient: chatroom.NewClient(conf),
		log:       log,
		rand: random{
			min: 10,
			max: random{30, 60}.get(),
		},
	}
}

func (bot *ChatterBot) Login(ctx context.Context, user *apiv1.UserDetails) error {
	if user == nil {
		user = &apiv1.UserDetails{
			Name: faker.FirstName(),
		}
		n := rand.IntN(9)
		if n < 6 {
			user.Name += " " + faker.FirstName()
		}
		if n < 3 {
			user.Name += " " + faker.LastName()
		}
	}

	return bot.botClient.Login(ctx, user)
}

func (bot *ChatterBot) Chatter(ctx context.Context) error {
	timer := time.NewTimer(time.Second * time.Duration(rand.IntN(bot.rand.min)))
	defer timer.Stop()

	var nextMsg string

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-timer.C:
			if nextMsg == "" {
				nextMsg = faker.Sentence()
				if _, err := bot.IndicateTyping(ctx, connect.NewRequest(&apiv1.IndicateTypingRequest{
					Typing: true,
				})); err != nil {
					return err
				}

				// fake typing of msg
				wait := bot.rand.duration(time.Millisecond * 20 * time.Duration(len(strings.Fields(nextMsg))))
				bot.log.Debug().Dur("wait", wait).Str("next", nextMsg).Msg("indicate typing")
				timer.Reset(wait)
			} else {
				if _, err := bot.SendChat(ctx, connect.NewRequest(&apiv1.SendChatRequest{
					Time: timestamppb.Now(),
					Text: nextMsg,
				})); err != nil {
					return err
				}

				wait := bot.rand.duration(time.Second)
				bot.log.Debug().Dur("wait", wait).Str("msg", nextMsg).Msg("send chat")
				timer.Reset(wait)
				nextMsg = ""
			}
		}
	}
}

type random struct {
	min int
	max int
}

func (r random) get() int { return r.min + rand.IntN(r.max-r.min) }

func (r random) duration(base time.Duration) time.Duration {
	return base * time.Duration(r.get())
}
