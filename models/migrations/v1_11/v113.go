// Copyright 2019 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_11

import (
	"fmt"

	"creepercoding.dev/models/db"
)

func FeatureChangeTargetBranch(x db.EngineMigration) error {
	type Comment struct {
		OldRef string
		NewRef string
	}

	if err := x.Sync(new(Comment)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}
	return nil
}
