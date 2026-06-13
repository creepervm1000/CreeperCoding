// Copyright 2019 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_10

import "creepercoding.dev/models/db"

func AddRepoAdminChangeTeamAccessColumnForUser(x db.EngineMigration) error {
	type User struct {
		RepoAdminChangeTeamAccess bool `xorm:"NOT NULL DEFAULT false"`
	}

	return x.Sync(new(User))
}
