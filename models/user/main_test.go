// Copyright 2021 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user_test

import (
	"testing"

	"creepercoding.dev/models/unittest"

	_ "creepercoding.dev/models"
	_ "creepercoding.dev/models/actions"
	_ "creepercoding.dev/models/activities"
	_ "creepercoding.dev/models/user"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
