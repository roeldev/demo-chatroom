// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { getContext, hasContext, setContext } from "svelte";
import { type Client, type Transport, createClient } from "@connectrpc/connect";
import { API } from "$lib/chatapi";
import { type EventStreamHandler } from "$lib/chatevents";
import { type User, mapUserDetails, type UUID, mapEventUser, type UserDetails } from "./user.ts";

const ctxKey = "users-list";

export function getUsersList(): UsersList {
    if (!hasContext(ctxKey)) {
        return setContext(ctxKey, new UsersList());
    }

    return getContext<UsersList>(ctxKey);
}

export class UsersList implements EventStreamHandler {
    private users: User[] = $state([]);
    private readonly client: Client<typeof API.RegistryService>;

    constructor(transport?: Transport) {
        this.client = createClient(
            API.RegistryService,
            transport ?? getContext("transport"),
        );
    }

    entries(): User[] {
        return this.users;
    }

    async fetch() {
        return this.client
            .activeUsers({})
            .then((response) => {
                this.users = response.users
                    .map((user) => {
                        return {
                            id: user.id!.value,
                            details: mapUserDetails(user.details!),
                            flags: user.flags!,
                            status: user.status!,
                        };
                    })
                    .sort(sortUsers);
            });
    }

    add(user: User) {
        this.users.push(user);
        this.users = this.users.toSorted(sortUsers);
    }

    remove(uuid: UUID) {
        const i = this.users.findIndex((u) => u.id === uuid);
        if (i < 0) {
            return;
        }

        this.users.splice(i, 1);
    }

    update(uuid: UUID, user: UserDetails) {
        const i = this.users.findIndex((u) => u.id === uuid);
        if (i < 0) {
            return;
        }

        const nameChanged = this.users[i].details.name != user.name;
        this.users[i].details = user;
        if (nameChanged) {
            this.users = this.users.toSorted(sortUsers);
        }
    }

    setStatus(uuid: UUID, status: API.UserStatus) {
        const i = this.users.findIndex((u) => u.id === uuid);
        if (i < 0) {
            return;
        }

        this.users[i].status = status;
    }


    handleEventStream(res: API.EventStreamResponse): boolean {
        // console.log("UsersList: handle event", res.event.case)
        switch (res.event.case) {
            case "userJoin":
                let user = mapEventUser(res.event.value.user!);
                user.flags = res.event.value.flags;

                this.add(user);
                break;

            case "userLeave":
                this.remove(res.event.value.user!.id!.value)
                break;

            case "userUpdate":
                this.update(
                    res.event.value.user!.id!.value,
                    mapUserDetails(res.event.value.user!.details!),
                );
                break;

            case "userStatus":
                this.setStatus(
                    res.event.value.user!.id!.value,
                    res.event.value.status!,
                );
                break;
        }

        return true;
    }
}

function sortUsers(a: User, b: User) {
    if (a.details.name > b.details.name) {
        return 1;
    }
    if (a.details.name < b.details.name) {
        return -1;
    }
    return 0;
}
