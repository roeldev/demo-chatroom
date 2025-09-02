// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { tick } from "svelte";
import { SvelteMap } from "svelte/reactivity";

import { API } from "$lib/chatapi";
import { activeUserID } from "$lib/chatauth";
import { type EventStreamHandler, EventType } from "$lib/chatevents";
import { type UUID, type UserDetails, mapUserDetails } from "$lib/chatusers";

import { ChatView } from "./chat-view.svelte.ts";

export const globalView = ""

let viewport: ChatViewport;

export function initChatViewport(): ChatViewport {
    viewport = new ChatViewport();
    return viewport;
}

export function getActiveReceiver(): UUID {
    return viewport.getReceiver();
}

export class ChatViewport implements EventStreamHandler {
    scrollElem?: HTMLElement;
    private scrollAtEnd: boolean = true;

    private receiver: UUID = $state(globalView);
    private readonly views = new SvelteMap<UUID, ChatView>();

    constructor() {
        this.views.set(globalView, new ChatView(globalView))
    }

    private getView(receiver: UUID): ChatView {
        if (!this.views.has(receiver)) {
            this.views.set(receiver, new ChatView(receiver))
        }
        return this.views.get(receiver)!
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

    async updateAfterStreamEvent() {
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

    handleEventStream(res: API.EventStreamResponse): boolean {
        let prevEventUser = "";
        switch (res.event.case) {
            case "userJoin": {
                prevEventUser = "";
                this.getView(globalView).addEvent({
                    case: res.event.case,
                    value: EventType.userJoin(res.time!, res.event.value),
                });
                break;
            }
            case "userLeave": {
                prevEventUser = "";

                const uid = res.event.value.user!.id!.value
                const event = {
                    case: res.event.case,
                    value: EventType.userLeave(res.time!, res.event.value),
                }
                this.views.forEach(view => {
                    // remove all references of user x is typing
                    view.setTyping(false, uid);
                    view.addEvent(event)
                });
                break;
            }
            case "userUpdate": {
                prevEventUser = "";
                const event = {
                    case: res.event.case,
                    value: EventType.userUpdate(res.time!, res.event.value),
                }
                this.views.forEach((view) => {
                    view.addEvent(event);
                });
                break;
            }
            case "userTyping": {
                const uid = res.event.value.user!.id!.value
                if (activeUserID() == uid) {
                    // no need to display that the current user is typing
                    return false;
                }

                const view = this.getView(asReceiver(res.event.value.receiverId))
                res.event.value.typing
                    ? view.setTyping(true, uid, mapUserDetails(res.event.value.user!.details!))
                    : view.setTyping(false, uid);

                return false;
            }

            case "chatSent": {
                const uid = res.event.value.user!.id!.value
                const view = this.getView(asReceiver(res.event.value.receiverId))

                if (activeUserID() == uid) {
                    prevEventUser = "";
                    view.addEvent({
                        case: "sentChat",
                        value: EventType.sentChat(res.time!, res.event.value),
                    });
                    break;
                }

                view.setTyping(false, uid);
                view.addEvent({
                    case: "receivedChat",
                    value: EventType.receivedChat(
                        res.time!,
                        res.event.value,
                        prevEventUser != uid,
                    ),
                });
                prevEventUser = uid;
                break
            }

            default:
                prevEventUser = "";
                break;
        }
        return true;
    }
}

function asReceiver(receiver?: API.UUID): UUID {
    return receiver?.value ?? globalView;
}
