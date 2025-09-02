// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { getContext } from "svelte";
import { type Client, type Transport, createClient } from "@connectrpc/connect";
import { API } from "$lib/chatapi";

export interface EventStreamHandler {
    handleEventStream(res: API.EventStreamResponse): boolean
}

// EventStreamer
export class EventStreamer {
    private readonly client: Client<typeof API.EventsService>;
    private readonly handlers: EventStreamHandler[] = [];
    private active: boolean = true;

    constructor(handlers: EventStreamHandler[], transport?: Transport) {
        this.handlers = handlers;
        this.client = createClient(
            API.EventsService,
            transport ?? getContext('transport'),
        );
    }

    async stream(callbackFn: () => void) {
        while (this.active) {
            for await (const res of this.client.eventStream({})) {
                this.handlers.every(h => h.handleEventStream(res));
                if (!!callbackFn) {
                    callbackFn();
                }
            }
        }
    }
}
