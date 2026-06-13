// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_19

import (
	"creepercoding.dev/models/db"
	"creepercoding.dev/modules/setting"
)

// AlterPublicGPGKeyImportContentFieldToMediumText: set GPGKeyImport Content field to MEDIUMTEXT
func AlterPublicGPGKeyImportContentFieldToMediumText(x db.EngineMigration) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	if setting.Database.Type.IsMySQL() {
		if _, err := sess.Exec("ALTER TABLE `gpg_key_import` CHANGE `content` `content` MEDIUMTEXT"); err != nil {
			return err
		}
	}
	return sess.Commit()
}
