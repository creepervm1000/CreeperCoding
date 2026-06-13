// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_20

import "creepercoding.dev/models/db"

func AddNewColumnForProject(x db.EngineMigration) error {
	type Project struct {
		OwnerID int64 `xorm:"INDEX"`
	}

	return x.Sync(new(Project))
}
