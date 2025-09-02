// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import * as EventType from "./event-types.ts";

export type Event = {
    case: "userJoin";
    value: EventType.UserJoinEvent;
} | {
    case: "userLeave";
    value: EventType.UserLeaveEvent;
} | {
    case: "userUpdate";
    value: EventType.UserUpdateEvent;
} | {
    case: "userStatus";
    value: EventType.UserStatusEvent;
} | {
    case: "receivedChat";
    value: EventType.ReceivedChatEvent;
} | {
    case: "sentChat";
    value: EventType.SentChatEvent;
} | { case: undefined; value?: undefined };
