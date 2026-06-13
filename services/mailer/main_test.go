// Copyright 2019 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package mailer

import (
	"testing"

	"creepercoding.dev/models/unittest"

	_ "creepercoding.dev/models"
	_ "creepercoding.dev/models/actions"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
