// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { getContext } from "svelte";
import { timestampNow } from "@bufbuild/protobuf/wkt";
import { createClient, type Client } from "@connectrpc/connect";

import { getActiveReceiver } from "$lib";
import { API } from "$lib/chatapi";

export class ChatInput {
    value: string = $state("")

    private readonly client: Client<typeof API.UserService>;
    private indicateTimestamp: number = 0

    constructor() {
        this.client = createClient(API.UserService, getContext("transport"));
    }

    indicateTyping(typing: boolean, timestamp?: number) {
        this.indicateTimestamp = typing ? timestamp! : 0;

        this.client
            .indicateTyping({
                receiverId: receiver(),
                typing: typing,
            })
            .catch((err) => {
                console.log("failed to indicate typing", err.message);
            })
    }

    onInput(event: Event) {
        if (this.value === "") {
            this.indicateTyping(false);
            return;
        }

        if ((event.timeStamp - this.indicateTimestamp) < 4000) {
            return;
        }
        this.indicateTyping(true, event.timeStamp);
    }

    onSubmitForm(event: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }) {
        event.preventDefault();
        this.indicateTyping(false);

        this.client
            .sendChat({
                receiverId: receiver(),
                time: timestampNow(),
                text: this.value,
            })
            .then(() => {
                this.value = "";
            })
            .catch((err) => {
                console.log("error", err.message);
            });
    }
}

function receiver() {
    const uuid = getActiveReceiver();
    return (uuid == "") ? undefined : { value: uuid };
}
