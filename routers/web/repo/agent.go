// Copyright 2025 The CreeperCoding Authors
// SPDX-License-Identifier: MIT

package repo

import (
	"net/http"

	"creepercoding.dev/modules/setting"
	"creepercoding.dev/modules/templates"
	"creepercoding.dev/services/context"
)

const tplAgent templates.TplName = "repo/agent"

func Agent(ctx *context.Context) {
	if !setting.Config().Ccopilot.Enabled.Value(ctx) {
		ctx.NotFound(nil)
		return
	}
	if ctx.Repo.Repository.CcopilotDisabled {
		ctx.NotFound(nil)
		return
	}
	ctx.Data["Title"] = ctx.Tr("repo.agent")
	ctx.Data["PageIsAgent"] = true
	ctx.HTML(http.StatusOK, tplAgent)
}
