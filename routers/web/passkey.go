// Copyright 2024 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package web

import (
	"net/http"

	"creepercoding.dev/modules/setting"
	"creepercoding.dev/services/context"
)

type passkeyEndpointsType struct {
	Enroll string `json:"enroll"`
	Manage string `json:"manage"`
}

func passkeyEndpoints(ctx *context.Context) {
	url := setting.AppURL + "user/settings/security"
	ctx.JSON(http.StatusOK, passkeyEndpointsType{
		Enroll: url,
		Manage: url,
	})
}
