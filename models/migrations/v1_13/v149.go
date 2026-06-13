// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_13

import (
	"fmt"

	"creepercoding.dev/models/db"
	"creepercoding.dev/modules/timeutil"
)

func AddCreatedAndUpdatedToMilestones(x db.EngineMigration) error {
	type Milestone struct {
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
	}

	if err := x.Sync(new(Milestone)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}
	return nil
}
