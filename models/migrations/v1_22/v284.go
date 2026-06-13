// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_22

import (
	"creepercoding.dev/models/db"

	"xorm.io/xorm"
)

func AddIgnoreStaleApprovalsColumnToProtectedBranchTable(x db.EngineMigration) error {
	type ProtectedBranch struct {
		IgnoreStaleApprovals bool `xorm:"NOT NULL DEFAULT false"`
	}
	_, err := x.SyncWithOptions(xorm.SyncOptions{
		IgnoreIndices:    true,
		IgnoreConstrains: true,
	}, new(ProtectedBranch))
	return err
}
