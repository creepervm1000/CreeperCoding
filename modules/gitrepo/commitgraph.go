// Copyright 2025 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package gitrepo

import (
	"context"

	"creepercoding.dev/modules/git"
)

func WriteCommitGraph(ctx context.Context, repo Repository) error {
	return git.WriteCommitGraph(ctx, repoPath(repo))
}
