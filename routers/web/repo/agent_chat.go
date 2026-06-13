// Copyright 2025 The CreeperCoding Authors
// SPDX-License-Identifier: MIT

package repo

import (
	"encoding/json"
	"net/http"

	"creepercoding.dev/modules/log"
	"creepercoding.dev/modules/setting"
	"creepercoding.dev/services/ccopilot"
	"creepercoding.dev/services/context"
)

func AgentChat(ctx *context.Context) {
	if !setting.Config().Ccopilot.Enabled.Value(ctx) {
		ctx.JSONError("ccopilot is not enabled")
		return
	}
	if ctx.Repo.Repository.CcopilotDisabled {
		ctx.JSONError("ccopilot is disabled for this repository")
		return
	}

	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(ctx.Req.Body).Decode(&req); err != nil {
		ctx.JSONError("invalid request")
		return
	}

	reply, err := ccopilot.AgentChat(ctx, ctx.Repo.Repository, ctx.Doer, req.Message)
	if err != nil {
		log.Error("ccopilot agent chat error: %v", err)
		ctx.JSONError("ccopilot error: " + err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"reply": reply})
}
