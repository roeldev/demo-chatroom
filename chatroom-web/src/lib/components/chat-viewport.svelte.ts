// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { tick } from "svelte";
import { SvelteMap } from "svelte/reactivity";

import { API } from "$lib/chatapi";
import { activeUserID } from "$lib/chatauth";
import { type EventStreamHandler, type PreviousEventHandler, EventType } from "$lib/chatevents";
import { type UUID, type UserDetails, mapUserDetails } from "$lib/chatusers";

import { ChatView } from "./chat-view.svelte.ts";
import type { UserJoinEvent } from "$api/v1/apiv1_pb.ts";

export const globalView = ""

let viewport: ChatViewport;

export function initChatViewport(): ChatViewport {
    viewport = new ChatViewport();
    return viewport;
}

export function getActiveReceiver(): UUID {
    return viewport.getReceiver();
}

export class ChatViewport {
    scrollElem?: HTMLElement;
    private scrollAtEnd: boolean = true;

    private receiver: UUID = $state(globalView);
    private readonly views = new SvelteMap<UUID, ChatView>();
    private readonly handler: EventHandler;

    constructor() {
        this.views.set(globalView, new ChatView(globalView))
        this.handler = new EventHandler(this);
    }

    getEventsHandler(): EventStreamHandler & PreviousEventHandler {
        return this.handler;
    }

    getReceiver(): UUID {
        return this.receiver
    }

    getViewsEntries() {
        return this.views.entries()
    }

    // Get a SvelteMap of the typing users of the active view.
    getUsersTyping() {
        if (!this.views.has(this.receiver)) {
            return new SvelteMap<UUID, UserDetails>()
        }
        return this.views.get(this.receiver)!.getUsersTyping();
    }

    activateGlobalView() {
        this.receiver = globalView;
    }

    activateUserView(uuid: UUID) {
        this.receiver = uuid;
    }

    onScrollEnd() {
        if (!this.scrollElem) {
            return;
        }

        this.scrollAtEnd = this.scrollElem.clientHeight + this.scrollElem.scrollTop >= this.scrollElem.scrollHeight;
    }

    async updateScrollPosition() {
        if (!this.scrollAtEnd) {
            // events view is not at bottom, user is probably reading a previous message
            return;
        }
        if (!this.scrollElem) {
            // no scroll elem known, skip auto update
            return;
        }

        await tick();
        this.scrollElem.scroll({
            top: this.scrollElem.scrollHeight,
            behavior: "smooth",
        });
    }
}

class EventHandler implements EventStreamHandler, PreviousEventHandler {
    private readonly owner: ChatViewport

    constructor(owner: ChatViewport) {
        this.owner = owner;
    }

    private getView(receiver: UUID): ChatView {
        if (!this.owner["views"].has(receiver)) {
            this.owner["views"].set(receiver, new ChatView(receiver))
        }
        return this.owner["views"].get(receiver)!
    }

    private allViews(callback: (view: ChatView) => void) {
        this.owner["views"].forEach(callback);
    }

    handlePreviousEvent(res: API.PreviousEventsResponse): boolean {
        res.history.forEach((history) => {
            switch (history.event.case) {
                case "userJoin": {
                    this.getView(globalView)
                        .prependEvent(EventType.userJoin(
                            history.time!,
                            history.event.value,
                        ));
                    break;
                }
                case "userLeave": {
                    const event = EventType.userLeave(
                        history.time!,
                        history.event.value,
                    );
                    this.allViews(view => view.prependEvent(event));
                    break;
                }
                case "userUpdate": {
                    const event = EventType.userUpdate(history.time!, history.event.value);
                    this.allViews((view) => view.prependEvent(event));
                    break;
                }
                case "chatSent": {
                    this.getView(asReceiver(history.event.value.receiverId))
                        .prependEvent(EventType.receivedChat(
                            history.time!,
                            history.event.value,
                        ));
                    break
                }
            }
        })
        return true;
    }

    handleEventStream(res: API.EventStreamResponse): boolean {
        switch (res.event.case) {
            case "userJoin": {
                this.getView(globalView)
                    .appendEvent(EventType.userJoin(res.time!, res.event.value));
                break;
            }
            case "userLeave": {
                const uid = res.event.value.user!.id!.value;
                const event = EventType.userLeave(res.time!, res.event.value);
                this.allViews(view => {
                    // remove all references of user x is typing
                    view.setTyping(false, uid);
                    view.appendEvent(event);
                });
                break;
            }
            case "userUpdate": {
                const event = EventType.userUpdate(res.time!, res.event.value);
                this.allViews((view) => view.appendEvent(event));
                break;
            }
            case "userTyping": {
                const uid = res.event.value.user!.id!.value;
                if (activeUserID() == uid) {
                    // no need to display that the current user is typing
                    return false;
                }

                const view = this.getView(asReceiver(res.event.value.receiverId));
                res.event.value.typing
                    ? view.setTyping(true, uid, mapUserDetails(res.event.value.user!.details!))
                    : view.setTyping(false, uid);

                return false;
            }
            case "chatSent": {
                const uid = res.event.value.user!.id!.value;
                const view = this.getView(asReceiver(res.event.value.receiverId));

                if (activeUserID() == uid) {
                    view.appendEvent(EventType.sentChat(res.time!, res.event.value));
                    break;
                }

                view.setTyping(false, uid);
                view.appendEvent(EventType.receivedChat(res.time!, res.event.value));
                break;
            }
        }
        return true;
    }
}

function asReceiver(receiver?: API.UUID): UUID {
    return receiver?.value ?? globalView;
}
