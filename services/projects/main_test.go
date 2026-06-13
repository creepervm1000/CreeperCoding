// Copyright 2025 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package project

import (
	"testing"

	"creepercoding.dev/models/unittest"

	_ "creepercoding.dev/models/actions"
	_ "creepercoding.dev/models/activities"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
