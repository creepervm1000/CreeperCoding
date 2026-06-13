// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_21

import "creepercoding.dev/models/db"

func AddVersionToActionRunTable(x db.EngineMigration) error {
	type ActionRun struct {
		Version int `xorm:"version default 0"`
	}
	return x.Sync(new(ActionRun))
}
