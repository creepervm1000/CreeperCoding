// Copyright 2024 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package markdown

import (
	"os"
	"testing"

	"creepercoding.dev/modules/markup"
	"creepercoding.dev/modules/setting"
)

func TestMain(m *testing.M) {
	setting.IsInTesting = true
	markup.RenderBehaviorForTesting.DisableAdditionalAttributes = true
	os.Exit(m.Run())
}
