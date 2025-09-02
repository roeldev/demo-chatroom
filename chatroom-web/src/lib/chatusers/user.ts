// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { API } from "$lib/chatapi";

export type UUID = string;

export type User = {
    id: UUID;
    details: UserDetails;
    flags: API.UserFlag;
    status: API.UserStatus;
    typing?: boolean;
};

export type UserDetails = {
    name: string;
    initials: string;
    color1: string;
    color2: string;
};

export function mapEventUser(user: API.EventUser): User {
    return {
        id: user.id!.value,
        details: mapUserDetails(user.details!),
        flags: API.UserFlag.NONE,
        status: API.UserStatus.DEFAULT,
    }
}

export function mapUserDetails(user: API.UserDetails): UserDetails {
    return {
        name: user.name,
        initials: user.initials,
        color1: user.color1!.value,
        color2: user.color2!.value,
    };
}
