// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { getContext } from "svelte";
import { type Client, type Transport, createClient } from "@connectrpc/connect";
import { API } from "$lib/chatapi";

export interface EventStreamHandler {
    handleEventStream(res: API.EventStreamResponse): boolean
}

export interface PreviousEventHandler {
    handlePreviousEvent(res: API.PreviousEventsResponse): boolean
}

// EventStreamer
export class EventStreamer {
    private readonly client: Client<typeof API.EventsService>;
    private active: boolean = true;

    constructor(transport?: Transport) {
        this.client = createClient(
            API.EventsService,
            transport ?? getContext('transport'),
        );
    }

    loadPrevious(handlers: PreviousEventHandler[]) {
        this.client.previousEvents({
            limit: 100,
        }).then((res) => {
            handlers.every(h => h.handlePreviousEvent(res));
        }).catch((err) => {
            console.log("err loading previous events", err)
        });
    }

    async stream(handlers: EventStreamHandler[], callbackFn: () => void) {
        while (this.active) {
            for await (const res of this.client.eventStream({})) {
                handlers.every(h => h.handleEventStream(res));
                if (!!callbackFn) {
                    callbackFn();
                }
            }
        }
    }
}
