// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package project

import (
	"testing"

	"creepercoding.dev/models/unittest"

	_ "creepercoding.dev/models/repo"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{
		FixtureFiles: []string{
			"project.yml",
			"project_board.yml",
			"project_issue.yml",
			"repository.yml",
		},
	})
}
