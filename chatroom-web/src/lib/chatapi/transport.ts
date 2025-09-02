// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import { setContext } from "svelte";
import { createConnectTransport } from "@connectrpc/connect-web";
import { tokenAuthorizer } from "$lib/chatauth";

export function initConnectTransport() {
    const env = {
        API_SERVER_SCHEME: "http",
        API_SERVER_HOST: "localhost",
        API_SERVER_PORT: 8080,
    }

    setContext("transport", createConnectTransport({
        baseUrl: env.API_SERVER_SCHEME + "://" + env.API_SERVER_HOST + ":" + env.API_SERVER_PORT + "/",
        interceptors: [tokenAuthorizer],
    }));
}
