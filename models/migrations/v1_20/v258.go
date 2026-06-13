// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_20

import "creepercoding.dev/models/db"

func AddPinOrderToIssue(x db.EngineMigration) error {
	type Issue struct {
		PinOrder int `xorm:"DEFAULT 0"`
	}

	return x.Sync(new(Issue))
}
