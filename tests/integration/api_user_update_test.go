// Copyright 2024 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"net/http"
	"testing"

	auth_model "creepercoding.dev/models/auth"
	"creepercoding.dev/tests"
)

func TestAPIUpdateUser(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	normalUsername := "user2"
	session := loginUser(t, normalUsername)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteUser)

	req := NewRequestWithJSON(t, "PATCH", "/api/v1/user/settings", map[string]string{
		"website": "https://gitea.com",
	}).AddTokenAuth(token)
	MakeRequest(t, req, http.StatusOK)
}
