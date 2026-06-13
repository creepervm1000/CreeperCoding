// Copyright 2021 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package oauth2

import (
	"context"

	user_model "creepercoding.dev/models/user"
	"creepercoding.dev/services/auth/source/db"
)

// Authenticate falls back to the db authenticator
func (source *Source) Authenticate(ctx context.Context, user *user_model.User, login, password string) (*user_model.User, error) {
	return db.Authenticate(ctx, user, login, password)
}

// NB: Oauth2 does not implement LocalTwoFASkipper for password authentication
// as its password authentication drops to db authentication
