// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package system_test

import (
	"testing"

	"creepercoding.dev/models/unittest"

	_ "creepercoding.dev/models" // register models
	_ "creepercoding.dev/models/actions"
	_ "creepercoding.dev/models/activities"
	_ "creepercoding.dev/models/system" // register models of system
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
