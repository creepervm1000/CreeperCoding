// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build windows

package private

import (
	"net/http"

	"creepercoding.dev/modules/graceful"
	"creepercoding.dev/modules/private"
	"creepercoding.dev/services/context"
)

// Restart is not implemented for Windows based servers as they can't fork
func Restart(ctx *context.PrivateContext) {
	ctx.JSON(http.StatusNotImplemented, private.Response{
		UserMsg: "windows servers cannot be gracefully restarted - shutdown and restart manually",
	})
}

// Shutdown causes the server to perform a graceful shutdown
func Shutdown(ctx *context.PrivateContext) {
	graceful.GetManager().DoGracefulShutdown()
	ctx.PlainText(http.StatusOK, "success")
}
