// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package issues_test

import (
	"testing"

	issues_model "creepercoding.dev/models/issues"
	"creepercoding.dev/models/unittest"

	_ "creepercoding.dev/models"
	_ "creepercoding.dev/models/actions"
	_ "creepercoding.dev/models/activities"
	_ "creepercoding.dev/models/repo"
	_ "creepercoding.dev/models/user"

	"github.com/stretchr/testify/assert"
)

func TestFixturesAreConsistent(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	unittest.CheckConsistencyFor(t,
		&issues_model.Issue{},
		&issues_model.PullRequest{},
		&issues_model.Milestone{},
		&issues_model.Label{},
	)
}

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
