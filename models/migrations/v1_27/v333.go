// Copyright 2026 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_27

import (
	"creepercoding.dev/models/db"

	"xorm.io/xorm"
)

func AddBranchProtectionBypassAllowlist(x db.EngineMigration) error {
	type ProtectedBranch struct {
		EnableBypassAllowlist  bool    `xorm:"NOT NULL DEFAULT false"`
		BypassAllowlistUserIDs []int64 `xorm:"JSON TEXT"`
		BypassAllowlistTeamIDs []int64 `xorm:"JSON TEXT"`
	}

	_, err := x.SyncWithOptions(xorm.SyncOptions{
		IgnoreConstrains: true,
		IgnoreIndices:    true,
	}, new(ProtectedBranch))
	return err
}
