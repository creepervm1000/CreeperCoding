// Copyright 2025 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package gitrepo

import (
	"context"

	"creepercoding.dev/modules/git"
)

func NewBatch(ctx context.Context, repo Repository) (git.CatFileBatchCloser, error) {
	return git.NewBatch(ctx, repoPath(repo))
}
