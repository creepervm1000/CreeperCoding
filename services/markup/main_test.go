// Copyright 2022 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package markup

import (
	"testing"

	"creepercoding.dev/models/unittest"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{
		FixtureFiles: []string{"user.yml", "repository.yml", "access.yml", "repo_unit.yml", "issue.yml"},
	})
}
