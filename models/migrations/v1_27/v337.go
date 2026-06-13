// Copyright 2026 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_27

import (
	"creepercoding.dev/models/db"

	"xorm.io/xorm"
)

type repositoryWithCcopilotDisabled struct {
	CcopilotDisabled bool `xorm:"NOT NULL DEFAULT false"`
}

func (repositoryWithCcopilotDisabled) TableName() string {
	return "repository"
}

func AddCcopilotDisabledToRepository(x db.EngineMigration) error {
	_, err := x.SyncWithOptions(xorm.SyncOptions{
		IgnoreDropIndices: true,
	}, new(repositoryWithCcopilotDisabled))
	return err
}
