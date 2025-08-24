// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"path/filepath"
	"runtime"

	"github.com/go-pogo/env/envfile"
	"github.com/go-pogo/errors"
	"github.com/roeldev/demo-chatroom"
	"github.com/roeldev/demo-chatroom/app/welcome-bot"
)

func main() {
	_, base, _, _ := runtime.Caller(0)
	base = filepath.Dir(base)

	for dir, conf := range map[string]any{
		"api-server":  chatroom.Config{},
		"welcome-bot": welcomebot.Config{},
	} {
		dir = filepath.Join(base, dir)
		errors.FatalOnErr(envfile.Generate(dir, filepath.Join(dir, ".env"), conf))
	}
}
