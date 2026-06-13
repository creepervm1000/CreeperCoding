// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_21

import "creepercoding.dev/models/db"

// UpdateActionsRefIndex updates the index of actions ref field
func UpdateActionsRefIndex(x db.EngineMigration) error {
	type ActionRun struct {
		Ref string `xorm:"index"` // the commit/tag/… causing the run
	}
	return x.Sync(new(ActionRun))
}
