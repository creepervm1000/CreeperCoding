// Copyright 2025 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package gitrepo

import (
	"context"

	"creepercoding.dev/modules/git"
)

func GetSigningKey(ctx context.Context) (*git.SigningKey, *git.Signature) {
	return git.GetSigningKey(ctx)
}
