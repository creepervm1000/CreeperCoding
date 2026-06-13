// Copyright 2022 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_17

import "creepercoding.dev/models/db"

func AddAllowMaintainerEdit(x db.EngineMigration) error {
	// PullRequest represents relation between pull request and repositories.
	type PullRequest struct {
		AllowMaintainerEdit bool `xorm:"NOT NULL DEFAULT false"`
	}

	return x.Sync(new(PullRequest))
}
