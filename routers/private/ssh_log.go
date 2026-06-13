// Copyright 2021 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package private

import (
	"net/http"

	"creepercoding.dev/modules/log"
	"creepercoding.dev/modules/private"
	"creepercoding.dev/modules/setting"
	"creepercoding.dev/modules/web"
	"creepercoding.dev/services/context"
)

// SSHLog hook to response ssh log
func SSHLog(ctx *context.PrivateContext) {
	if !setting.Log.EnableSSHLog {
		ctx.Status(http.StatusOK)
		return
	}

	opts := web.GetForm(ctx).(*private.SSHLogOption)

	if opts.IsError {
		log.Error("ssh: %v", opts.Message)
		ctx.Status(http.StatusOK)
		return
	}

	log.Debug("ssh: %v", opts.Message)
	ctx.Status(http.StatusOK)
}
