// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_20

import "creepercoding.dev/models/db"

func AddVersionToActionRunner(x db.EngineMigration) error {
	type ActionRunner struct {
		Version string `xorm:"VARCHAR(64)"` // the version of act_runner
	}

	return x.Sync(new(ActionRunner))
}
