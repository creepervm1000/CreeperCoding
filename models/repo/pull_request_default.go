// Copyright 2026 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"context"

	"creepercoding.dev/models/unit"
	"creepercoding.dev/modules/util"
)

func (repo *Repository) GetPullRequestTargetBranch(ctx context.Context) string {
	unitPRConfig := repo.MustGetUnit(ctx, unit.TypePullRequests).PullRequestsConfig()
	return util.IfZero(unitPRConfig.DefaultTargetBranch, repo.DefaultBranch)
}
