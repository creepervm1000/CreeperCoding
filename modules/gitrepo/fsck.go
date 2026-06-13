// Copyright 2024 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package gitrepo

import (
	"context"
	"time"

	"creepercoding.dev/modules/git/gitcmd"
)

// Fsck verifies the connectivity and validity of the objects in the database
func Fsck(ctx context.Context, repo Repository, timeout time.Duration, args gitcmd.TrustedCmdArgs) error {
	return RunCmd(ctx, repo, gitcmd.NewCommand("fsck").AddArguments(args...).WithTimeout(timeout))
}
