// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1connect

import (
	"context"
	"net/http"
	"sync"

	"connectrpc.com/connect"
	apiv1 "github.com/roeldev/demo-chatroom/api/v1"
)

// https://connectrpc.com/docs/go/interceptors

type ClientInterceptor struct {
	mut   sync.RWMutex
	token string
}

func NewClientInterceptor(token string) *ClientInterceptor {
	return &ClientInterceptor{token: token}
}

func (ci *ClientInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if !req.Spec().IsClient {
			panic("apiv1connect: Use NewHandlerInterceptor() instead")
		}

		ci.addAuthHeader(req.Header())
		res, err := next(ctx, req)
		if err != nil {
			return res, err
		}

		if join, ok := res.Any().(*apiv1.JoinResponse); ok {
			ci.UpdateToken(join.Token)
		}
		return res, nil
	}
}

func (ci *ClientInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		conn := next(ctx, spec)
		ci.addAuthHeader(conn.RequestHeader())
		return conn
	}
}

func (ci *ClientInterceptor) addAuthHeader(h http.Header) {
	ci.mut.RLock()
	defer ci.mut.RUnlock()

	if ci.token != "" {
		h.Set("authorization", "bearer "+ci.token)
	}
}

func (ci *ClientInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

func (ci *ClientInterceptor) UpdateToken(token string) {
	ci.mut.Lock()
	ci.token = token
	ci.mut.Unlock()
}
