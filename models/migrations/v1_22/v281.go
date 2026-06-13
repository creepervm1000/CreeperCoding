// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_22

import (
	"creepercoding.dev/models/db"
	"creepercoding.dev/modules/timeutil"
)

func CreateAuthTokenTable(x db.EngineMigration) error {
	type AuthToken struct {
		ID          string `xorm:"pk"`
		TokenHash   string
		UserID      int64              `xorm:"INDEX"`
		ExpiresUnix timeutil.TimeStamp `xorm:"INDEX"`
	}

	return x.Sync(new(AuthToken))
}
