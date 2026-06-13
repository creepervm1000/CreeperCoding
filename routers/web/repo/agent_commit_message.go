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

func AgentCommitMessage(ctx *context.Context) {
	if !setting.Config().Ccopilot.Enabled.Value(ctx) {
		ctx.JSONError("ccopilot is not enabled")
		return
	}
	if ctx.Repo.Repository.CcopilotDisabled {
		ctx.JSONError("ccopilot is disabled for this repository")
		return
	}

	var req struct {
		Path           string `json:"path"`
		Content        string `json:"content"`
		Branch         string `json:"branch"`
		CommitSummary  string `json:"commit_summary"`
	}
	if err := json.NewDecoder(ctx.Req.Body).Decode(&req); err != nil {
		ctx.JSONError("invalid request")
		return
	}

	msg, err := ccopilot.GenerateCommitMessage(ctx, ctx.Repo.Repository, req.Path, req.Content, req.Branch, req.CommitSummary)
	if err != nil {
		log.Error("ccopilot commit message error: %v", err)
		ctx.JSONError("ccopilot error: " + err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"message": msg})
}
