// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatbot

import (
	"github.com/go-pogo/webapp/logger"
	chatroom "github.com/roeldev/demo-chatroom"
)

type Config struct {
	Logger logger.Config         `env:",include"`
	Client chatroom.ClientConfig `env:",include"`
}

type botClient = chatroom.Client
