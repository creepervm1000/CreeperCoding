// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package convert

import (
	"testing"

	"creepercoding.dev/models/unittest"
	user_model "creepercoding.dev/models/user"
	api "creepercoding.dev/modules/structs"

	"github.com/stretchr/testify/assert"
)

func TestUser_ToUser(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	user1 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1, IsAdmin: true})

	apiUser := toUser(t.Context(), user1, true, true)
	assert.True(t, apiUser.IsAdmin)
	assert.Contains(t, apiUser.AvatarURL, "://")

	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2, IsAdmin: false})

	apiUser = toUser(t.Context(), user2, true, true)
	assert.False(t, apiUser.IsAdmin)

	apiUser = toUser(t.Context(), user1, false, false)
	assert.False(t, apiUser.IsAdmin)
	assert.Equal(t, api.UserVisibilityPublic, apiUser.Visibility)

	user31 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 31, IsAdmin: false, Visibility: api.VisibleTypePrivate})

	apiUser = toUser(t.Context(), user31, true, true)
	assert.False(t, apiUser.IsAdmin)
	assert.Equal(t, api.UserVisibilityPrivate, apiUser.Visibility)
}
