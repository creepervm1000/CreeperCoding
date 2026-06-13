// Copyright 2026 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package pipeline

import (
	"testing"

	"creepercoding.dev/modules/git"
)

func TestMain(m *testing.M) {
	git.RunGitTests(m)
}
