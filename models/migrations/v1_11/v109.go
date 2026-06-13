// Copyright 2019 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_11

import "creepercoding.dev/models/db"

func AddCanCreateOrgRepoColumnForTeam(x db.EngineMigration) error {
	type Team struct {
		CanCreateOrgRepo bool `xorm:"NOT NULL DEFAULT false"`
	}

	return x.Sync(new(Team))
}
