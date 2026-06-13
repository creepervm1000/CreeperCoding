// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user_test

import (
	"testing"

	"creepercoding.dev/models/unittest"
	user_model "creepercoding.dev/models/user"

	"github.com/stretchr/testify/assert"
)

func TestLookupUserRedirect(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	userID, err := user_model.LookupUserRedirect(t.Context(), "olduser1")
	assert.NoError(t, err)
	assert.EqualValues(t, 1, userID)

	_, err = user_model.LookupUserRedirect(t.Context(), "doesnotexist")
	assert.True(t, user_model.IsErrUserRedirectNotExist(err))
}
