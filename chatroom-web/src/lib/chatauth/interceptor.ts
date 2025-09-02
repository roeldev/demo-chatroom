// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

import type { Interceptor } from "@connectrpc/connect";
import { getAuthToken } from "./auth.ts";

// tokenAuthorizer is an interceptor which adds an Authorization header
// containing the auth token which was given to us when joining the chatroom.
export const tokenAuthorizer: Interceptor = (next) => async (req) => {
    const token = getAuthToken()
    if (token) {
        req.header.set("Authorization", "Bearer " + token);
    }
    return await next(req);
};
