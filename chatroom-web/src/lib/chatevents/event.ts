// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { type Timestamp, timestampDate } from "@bufbuild/protobuf/wkt";
import { API } from "$lib/chatapi";

export type Event = UserJoin | UserLeave | UserUpdate | UserStatus |
    ReceivedChat | SentChat | { case: undefined; value?: undefined; };

////////////////////////////////////////////////////////////////////////////////

type UserJoin = {
    case: "userJoin";
    value: UserJoinEvent;
}

export type UserJoinEvent = API.UserJoinEvent & {
    date: Date;
    user: API.EventUser;
}

export function userJoin(ts: Timestamp, event: API.UserJoinEvent): Event & UserJoin {
    return {
        case: "userJoin",
        value: {
            date: timestampDate(ts),
            user: event.user!,
            ...event
        },
    };
}

////////////////////////////////////////////////////////////////////////////////

type UserLeave = {
    case: "userLeave";
    value: UserLeaveEvent;
}

export type UserLeaveEvent = API.UserLeaveEvent & {
    date: Date;
    user: API.EventUser;
}

export function userLeave(ts: Timestamp, event: API.UserLeaveEvent): Event & UserLeave {
    return {
        case: "userLeave",
        value: {
            date: timestampDate(ts),
            user: event.user!,
            ...event
        },
    }
}

////////////////////////////////////////////////////////////////////////////////

type UserUpdate = {
    case: "userUpdate";
    value: UserUpdateEvent;
}

export type UserUpdateEvent = API.UserUpdateEvent & {
    date: Date;
    user: API.EventUser;
}

export function userUpdate(ts: Timestamp, event: API.UserUpdateEvent): Event & UserUpdate {
    return {
        case: "userUpdate",
        value: {
            date: timestampDate(ts),
            user: event.user!,
            ...event
        },
    }
}

////////////////////////////////////////////////////////////////////////////////

type UserStatus = {
    case: "userStatus";
    value: UserStatusEvent;
}

export type UserStatusEvent = API.UserStatusEvent & {
    date: Date;
    user: API.EventUser;
}

export function userStatus(ts: Timestamp, event: API.UserStatusEvent): Event & UserStatus {
    return {
        case: "userStatus",
        value: {
            date: timestampDate(ts),
            user: event.user!,
            ...event
        },
    }
}

////////////////////////////////////////////////////////////////////////////////

type ReceivedChat = {
    case: "receivedChat";
    value: ReceivedChatEvent;
}

export type ReceivedChatEvent = API.ChatSentEvent & {
    date: Date;
    showAvatar?: boolean;
    chatId: API.UUID;
    user: API.EventUser;
}

export function receivedChat(ts: Timestamp, event: API.ChatSentEvent): Event & ReceivedChat {
    return {
        case: "receivedChat",
        value: {
            date: timestampDate(ts),
            chatId: event.chatId!,
            user: event.user!,
            ...event
        },
    }
}

////////////////////////////////////////////////////////////////////////////////

type SentChat = {
    case: "sentChat";
    value: SentChatEvent;
}

export type SentChatEvent = API.ChatSentEvent & {
    date: Date;
    chatId: API.UUID;
}

export function sentChat(ts: Timestamp, event: API.ChatSentEvent): Event & SentChat {
    return {
        case: "sentChat",
        value: {
            date: timestampDate(ts),
            chatId: event.chatId!,
            ...event
        },
    }
}
