// Copyright 2021 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_16

import (
	"creepercoding.dev/models/db"
	"creepercoding.dev/models/migrations/base"
)

func RecreateIssueResourceIndexTable(x db.EngineMigration) error {
	type IssueIndex struct {
		GroupID  int64 `xorm:"pk"`
		MaxIndex int64 `xorm:"index"`
	}

	return base.RecreateTables(new(IssueIndex))(x)
}
