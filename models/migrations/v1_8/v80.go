// Copyright 2019 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_8

import "creepercoding.dev/models/db"

func AddIsLockedToIssues(x db.EngineMigration) error {
	// Issue see models/issue.go
	type Issue struct {
		ID       int64 `xorm:"pk autoincr"`
		IsLocked bool  `xorm:"NOT NULL DEFAULT false"`
	}

	return x.Sync(new(Issue))
}
