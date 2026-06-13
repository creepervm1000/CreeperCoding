// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package db_test

import (
	"testing"

	"creepercoding.dev/models/unittest"

	_ "creepercoding.dev/models"
	_ "creepercoding.dev/models/repo"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
