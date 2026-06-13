// Copyright 2025 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_25

import (
	"creepercoding.dev/models/db"
	"creepercoding.dev/models/migrations/base"

	"xorm.io/xorm/schemas"
)

func ExtendCommentTreePathLength(x db.EngineMigration) error {
	dbType := x.Dialect().URI().DBType
	if dbType == schemas.SQLITE { // For SQLITE, varchar or char will always be represented as TEXT
		return nil
	}

	return base.ModifyColumn(x, "comment", &schemas.Column{
		Name: "tree_path",
		SQLType: schemas.SQLType{
			Name: "VARCHAR",
		},
		Length:         4000,
		Nullable:       true, // To keep compatible as nullable
		DefaultIsEmpty: true,
	})
}
