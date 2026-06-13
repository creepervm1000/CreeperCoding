// Copyright 2017 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package misc

import (
	"net/http"

	"creepercoding.dev/services/context"
)

func Swagger(ctx *context.Context) {
	ctx.HTML(http.StatusOK, "swagger/openapi-viewer")
}
