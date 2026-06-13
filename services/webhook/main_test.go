// Copyright 2019 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"testing"

	"creepercoding.dev/models/unittest"
	"creepercoding.dev/modules/hostmatcher"
	"creepercoding.dev/modules/setting"

	_ "creepercoding.dev/models"
	_ "creepercoding.dev/models/actions"
)

func TestMain(m *testing.M) {
	// for tests, allow only loopback IPs
	setting.Webhook.AllowedHostList = hostmatcher.MatchBuiltinLoopback
	unittest.MainTest(m, &unittest.TestOptions{
		SetUp: func() error {
			setting.LoadQueueSettings()
			return Init()
		},
	})
}
