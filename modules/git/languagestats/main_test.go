// Copyright 2025 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package languagestats

import (
	"testing"

	"creepercoding.dev/modules/git"
)

func TestMain(m *testing.M) {
	git.RunGitTests(m)
}
