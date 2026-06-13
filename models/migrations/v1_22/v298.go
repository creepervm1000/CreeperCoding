// Copyright 2024 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_22

import "creepercoding.dev/models/db"

func DropWronglyCreatedTable(x db.EngineMigration) error {
	return x.DropTables("o_auth2_application")
}
