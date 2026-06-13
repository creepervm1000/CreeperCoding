// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo_test

import (
	"testing"

	"creepercoding.dev/models/unittest"

	_ "creepercoding.dev/models" // register table model
	_ "creepercoding.dev/models/actions"
	_ "creepercoding.dev/models/activities"
	_ "creepercoding.dev/models/perm/access" // register table model
	_ "creepercoding.dev/models/repo"        // register table model
	_ "creepercoding.dev/models/user"        // register table model
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
