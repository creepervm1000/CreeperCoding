// Copyright 2025 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package admin

import (
	"net/http"

	"creepercoding.dev/modules/tailmsg"
	"creepercoding.dev/services/context"
)

func PerfTrace(ctx *context.Context) {
	monitorTraceCommon(ctx)
	ctx.Data["PageIsAdminMonitorPerfTrace"] = true
	ctx.Data["PerfTraceRecords"] = tailmsg.GetManager().GetTraceRecorder().GetRecords()
	ctx.HTML(http.StatusOK, tplPerfTrace)
}
