// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { type Timestamp, timestampDate } from "@bufbuild/protobuf/wkt";
import { API } from "$lib/chatapi";

export type UserJoinEvent = API.UserJoinEvent & {
    date: Date,
    user: API.EventUser;
}

export function userJoin(ts: Timestamp, event: API.UserJoinEvent): UserJoinEvent {
    return {
        date: timestampDate(ts),
        user: event.user!,
        ...event
    }
}

export type UserLeaveEvent = API.UserLeaveEvent & {
    date: Date,
    user: API.EventUser;
}

export function userLeave(ts: Timestamp, event: API.UserLeaveEvent): UserLeaveEvent {
    return {
        date: timestampDate(ts),
        user: event.user!,
        ...event
    }
}

export type UserUpdateEvent = API.UserUpdateEvent & {
    date: Date,
    user: API.EventUser;
}

export function userUpdate(ts: Timestamp, event: API.UserUpdateEvent): UserUpdateEvent {
    return {
        date: timestampDate(ts),
        user: event.user!,
        ...event
    }
}

export type UserStatusEvent = API.UserStatusEvent & {
    date: Date,
    user: API.EventUser,
}

export function userStatus(ts: Timestamp, event: API.UserStatusEvent): UserStatusEvent {
    return {
        date: timestampDate(ts),
        user: event.user!,
        ...event
    }
}

export type ReceivedChatEvent = API.ChatSentEvent & {
    date: Date,
    showAvatar?: boolean,
    chatId: API.UUID,
    user: API.EventUser,
}

export function receivedChat(ts: Timestamp, event: API.ChatSentEvent, showAvatar: boolean): ReceivedChatEvent {
    return {
        date: timestampDate(ts),
        showAvatar: showAvatar,
        chatId: event.chatId!,
        user: event.user!,
        ...event
    }
}

export type SentChatEvent = API.ChatSentEvent & {
    date: Date,
    chatId: API.UUID,
}

export function sentChat(ts: Timestamp, event: API.ChatSentEvent): SentChatEvent {
    return {
        date: timestampDate(ts),
        chatId: event.chatId!,
        ...event
    }
}
