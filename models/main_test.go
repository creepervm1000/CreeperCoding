// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package models

import (
	"testing"

	activities_model "creepercoding.dev/models/activities"
	"creepercoding.dev/models/organization"
	repo_model "creepercoding.dev/models/repo"
	"creepercoding.dev/models/unittest"
	user_model "creepercoding.dev/models/user"

	_ "creepercoding.dev/models/actions"
	_ "creepercoding.dev/models/system"

	"github.com/stretchr/testify/assert"
)

// TestFixturesAreConsistent assert that test fixtures are consistent
func TestFixturesAreConsistent(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	unittest.CheckConsistencyFor(t,
		&user_model.User{},
		&repo_model.Repository{},
		&organization.Team{},
		&activities_model.Action{})
}

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
