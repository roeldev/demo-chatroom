// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { SvelteMap } from "svelte/reactivity";

import { type Event } from '$lib/chatevents';
import { type UUID, type UserDetails } from "$lib/chatusers";

export class ChatView {
    readonly owner: UUID;
    private events: Event[] = $state([]);
    private usersTyping = new SvelteMap<UUID, UserDetails>();

    constructor(owner: UUID) {
        this.owner = owner;
    }

    appendEvent(e: Event): void {
        this.events.push(e);
    }

    prependEvent(e: Event): void {
        this.events.unshift(e);
    }

    getEvents() {
        return this.events
    }

    setTyping(typing: boolean, uid: UUID, user?: UserDetails) {
        typing ? this.usersTyping.set(uid, user!) : this.usersTyping.delete(uid)
    }

    getUsersTyping() {
        return this.usersTyping;
    }
}
