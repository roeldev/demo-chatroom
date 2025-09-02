// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { type JwtPayload, jwtDecode } from "jwt-decode";
import { get, type Writable, writable } from 'svelte/store';
import type { UUID } from "$lib/chatusers";

export type AuthClaims = JwtPayload & {
    uid?: UUID;
}

const authClaims: Writable<AuthClaims> = writable({});

// authToken store which is accessible throughout the application's life.
const authToken = writable("");

export function isAuthenticated(): boolean {
    return get(authToken) != "";
}

export function setAuthToken(token: string) {
    authToken.set(token);
    authClaims.set(jwtDecode(token))
}

export function getAuthToken(): string {
    return get(authToken);
}

export function getAuthClaims(): AuthClaims {
    return get(authClaims);
}

export function activeUserID(): string {
    return getAuthClaims().uid!
}
