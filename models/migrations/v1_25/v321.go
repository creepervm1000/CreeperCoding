// Copyright 2025 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_25

import (
	"creepercoding.dev/models/db"
	"creepercoding.dev/models/migrations/base"
	"creepercoding.dev/modules/setting"

	"xorm.io/xorm/schemas"
)

func UseLongTextInSomeColumnsAndFixBugs(x db.EngineMigration) error {
	if !setting.Database.Type.IsMySQL() {
		return nil // Only mysql need to change from text to long text, for other databases, they are the same
	}

	if err := base.ModifyColumn(x, "review_state", &schemas.Column{
		Name: "updated_files",
		SQLType: schemas.SQLType{
			Name: "LONGTEXT",
		},
		Length:         0,
		Nullable:       false,
		DefaultIsEmpty: true,
	}); err != nil {
		return err
	}

	if err := base.ModifyColumn(x, "package_property", &schemas.Column{
		Name: "value",
		SQLType: schemas.SQLType{
			Name: "LONGTEXT",
		},
		Length:         0,
		Nullable:       false,
		DefaultIsEmpty: true,
	}); err != nil {
		return err
	}

	return base.ModifyColumn(x, "notice", &schemas.Column{
		Name: "description",
		SQLType: schemas.SQLType{
			Name: "LONGTEXT",
		},
		Length:         0,
		Nullable:       false,
		DefaultIsEmpty: true,
	})
}
