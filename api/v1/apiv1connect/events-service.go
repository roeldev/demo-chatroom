// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package apiv1connect

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/go-pogo/errors"
	"github.com/google/uuid"
	"github.com/roeldev/demo-chatroom/api/v1"
	"github.com/roeldev/demo-chatroom/chatauth"
	"github.com/roeldev/demo-chatroom/chatevents"
	"github.com/roeldev/demo-chatroom/chatevents/event"
	"github.com/roeldev/demo-chatroom/chatusers"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ EventsServiceHandler = (*EventsService)(nil)

type EventsService struct {
	log     zerolog.Logger
	leaver  chatauth.Leaver
	history chatevents.EventsLister
	events  *eventHandler
}

func NewEventsService(log zerolog.Logger, history chatevents.EventsLister, broker *chatevents.EventsBroker, leaver chatauth.Leaver) *EventsService {
	svc := &EventsService{
		log:     log,
		leaver:  leaver,
		history: history,
		events:  newEventHandler(),
	}
	broker.Handle(svc.events)
	return svc
}

func (svc *EventsService) PreviousEvents(_ context.Context, req *connect.Request[apiv1.PreviousEventsRequest]) (*connect.Response[apiv1.PreviousEventsResponse], error) {
	var until time.Time
	if req.Msg.UntilTime.IsValid() {
		until = req.Msg.UntilTime.AsTime()
	}

	events := svc.history.ListEvents(until, req.Msg.Limit)
	history := make([]*apiv1.PreviousEventsResponse_PreviousEvent, 0, len(events))

	for _, evt := range events {
		history = append(history, &apiv1.PreviousEventsResponse_PreviousEvent{
			Time:  timestamppb.New(evt.Time),
			Event: apiv1.NewPreviousEventsResponseEvent(evt.Type),
		})
	}

	return connect.NewResponse(&apiv1.PreviousEventsResponse{
		History: history,
	}), nil
}

// EventStream streams chat related events to any connected client.
// https://connectrpc.com/docs/go/streaming
func (svc *EventsService) EventStream(ctx context.Context, _ *connect.Request[emptypb.Empty], stream *connect.ServerStream[apiv1.EventStreamResponse]) error {
	streamStart := time.Now()

	user := getUser(ctx)
	svc.log.Debug().
		Stringer("user", user).
		Msg("start stream")

	ch := svc.events.subscribe(user.ID)
	defer func() {
		svc.events.unsubscribe(user.ID)
		if svc.leaver != nil {
			svc.leaver.Leave(user.ID)
		}

		svc.log.Debug().
			Stringer("user", user).
			Stringer("duration", time.Since(streamStart)).
			Msg("close stream")
	}()

streamLoop:
	for {
		select {
		case evt, more := <-ch:
			if !more {
				break streamLoop
			}

			if ut, ok := evt.Type.(*event.UserTypingEvent); ok && ut.UserID == user.ID {
				// ignore current user typing events
				continue streamLoop
			}

			streamErr := stream.Send(&apiv1.EventStreamResponse{
				Time:  timestamppb.New(evt.Time),
				Event: apiv1.NewEventStreamResponseEvent(evt.Type),
			})
			if streamErr == nil {
				svc.log.Debug().
					Bool("more", more).
					Stringer("open", evt.Time.Sub(streamStart)).
					EmbedObject(evt).
					Stringer("to", user).
					Msg("sent stream event")
				continue streamLoop
			}

			// todo: return on connect.CodeUnauthenticated
			svc.log.Warn().
				Bool("more", more).
				Stringer("open", evt.Time.Sub(streamStart)).
				EmbedObject(evt).
				Stringer("user", user).
				Err(streamErr).
				Msg("stream error")

			var connectErr *connect.Error
			if errors.As(streamErr, &connectErr) &&
				connectErr.Code() == connect.CodeCanceled ||
				connectErr.Code() == connect.CodeDeadlineExceeded {

				// close stream
				break streamLoop
			}

		//case <-time.After(time.Until(claims.ExpiresAt.Time)):
		//	svc.log.Debug().Msg("token expired")
		//	break streamLoop

		case <-ctx.Done():
			break streamLoop
		}
	}
	return nil
}

type eventChan chan chatevents.Event

type eventHandler struct {
	mut  sync.RWMutex
	subs map[chatusers.UserID]eventChan
}

func newEventHandler() *eventHandler {
	return &eventHandler{
		subs: make(map[chatusers.UserID]eventChan),
	}
}

func (eh *eventHandler) subscribe(uid chatusers.UserID) eventChan {
	eh.mut.Lock()
	defer eh.mut.Unlock()

	// close existing channel
	if sub, ok := eh.subs[uid]; ok {
		close(sub)
	}

	sub := make(chan chatevents.Event, 10)
	eh.subs[uid] = sub
	return sub
}

func (eh *eventHandler) unsubscribe(uid chatusers.UserID) {
	eh.mut.Lock()
	defer eh.mut.Unlock()

	if sub, ok := eh.subs[uid]; ok {
		delete(eh.subs, uid)
		close(sub)
	}
}

func (eh *eventHandler) HandleEvent(e chatevents.Event) {
	eh.mut.RLock()
	defer eh.mut.RUnlock()

	for user, ch := range eh.subs {
		if re := e.AsReceiverEvent(); re != nil {
			receiver := re.GetReceiverID()
			if receiver != uuid.Nil && receiver != user {
				// event is not meant for user
				continue
			}
		}

		ch <- e
	}
}
