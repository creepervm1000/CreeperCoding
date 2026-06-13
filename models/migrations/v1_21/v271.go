// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_21

import (
	"creepercoding.dev/models/db"
	"creepercoding.dev/modules/timeutil"
)

func AddArchivedUnixColumInLabelTable(x db.EngineMigration) error {
	type Label struct {
		ArchivedUnix timeutil.TimeStamp `xorm:"DEFAULT NULL"`
	}
	return x.Sync(new(Label))
}
