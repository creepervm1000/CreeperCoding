// Copyright 2019 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_10

import "creepercoding.dev/models/db"

func ChangeSomeColumnsLengthOfExternalLoginUser(x db.EngineMigration) error {
	type ExternalLoginUser struct {
		AccessToken       string `xorm:"TEXT"`
		AccessTokenSecret string `xorm:"TEXT"`
		RefreshToken      string `xorm:"TEXT"`
	}

	return x.Sync(new(ExternalLoginUser))
}
