// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"testing"

	"creepercoding.dev/models/unittest"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{
		FixtureFiles: []string{
			"webhook.yml",
			"hook_task.yml",
		},
		SetUp: prepareWebhookTestData,
	})
}
