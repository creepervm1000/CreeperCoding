// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_21

import (
	"creepercoding.dev/models/db"
	"creepercoding.dev/modules/timeutil"
)

func CreateActionTasksVersionTable(x db.EngineMigration) error {
	type ActionTasksVersion struct {
		ID          int64 `xorm:"pk autoincr"`
		OwnerID     int64 `xorm:"UNIQUE(owner_repo)"`
		RepoID      int64 `xorm:"INDEX UNIQUE(owner_repo)"`
		Version     int64
		CreatedUnix timeutil.TimeStamp `xorm:"created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"updated"`
	}

	return x.Sync(new(ActionTasksVersion))
}
