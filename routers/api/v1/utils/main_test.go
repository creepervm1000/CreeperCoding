// Copyright 2018 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package utils

import (
	"testing"

	"creepercoding.dev/models/unittest"
	"creepercoding.dev/modules/setting"
	webhook_service "creepercoding.dev/services/webhook"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{
		SetUp: func() error {
			setting.LoadQueueSettings()
			return webhook_service.Init()
		},
	})
}
