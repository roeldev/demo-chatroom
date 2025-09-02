// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { getContext } from "svelte";
import { type Client, type Transport, createClient } from "@connectrpc/connect";

import { goto } from "$app/navigation";
import { API } from "$lib/chatapi";
import { setAuthToken } from "$lib/chatauth/auth.ts";

export type JoinSubmitEvent = SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }

export class JoinForm {
    private readonly client: Client<typeof API.AuthService>;
    private readonly fieldName: string

    constructor(fieldName: string, transport?: Transport) {
        this.fieldName = fieldName;
        this.client = createClient(
            API.AuthService,
            transport ?? getContext("transport"),
        );
    }

    async submitForm(event: JoinSubmitEvent) {
        event.preventDefault();
        const form = new FormData(event.currentTarget);
        const input = form.get(this.fieldName)?.toString()!;

        return this.client
            .join({
                user: {
                    name: input,
                }
            })
            .then((response) => {
                setAuthToken(response.token)
                goto("/")
            });
    }
}
