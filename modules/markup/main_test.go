// Copyright 2024 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package markup_test

import (
	"os"
	"testing"

	"creepercoding.dev/modules/markup"
	"creepercoding.dev/modules/setting"
)

func TestMain(m *testing.M) {
	setting.IsInTesting = true
	markup.RenderBehaviorForTesting.DisableAdditionalAttributes = true
	setting.Markdown.FileNamePatterns = []string{"*.md"}
	markup.RefreshFileNamePatterns()
	os.Exit(m.Run())
}
