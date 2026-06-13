// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_21

import "creepercoding.dev/models/db"

func AddScheduleIDForActionRun(x db.EngineMigration) error {
	type ActionRun struct {
		ScheduleID int64
	}
	return x.Sync(new(ActionRun))
}
