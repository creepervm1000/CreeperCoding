// Copyright 2022 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"net/http"

	"creepercoding.dev/modules/web"
	"creepercoding.dev/routers/api/actions/ping"
	"creepercoding.dev/routers/api/actions/runner"
)

func Routes(prefix string) *web.Router {
	m := web.NewRouter()

	path, handler := ping.NewPingServiceHandler()
	m.Post(path+"*", http.StripPrefix(prefix, handler).ServeHTTP)

	path, handler = runner.NewRunnerServiceHandler()
	m.Post(path+"*", http.StripPrefix(prefix, handler).ServeHTTP)

	return m
}
