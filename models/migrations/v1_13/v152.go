// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_13

import "creepercoding.dev/models/db"

func AddTrustModelToRepository(x db.EngineMigration) error {
	type Repository struct {
		TrustModel int
	}
	return x.Sync(new(Repository))
}
