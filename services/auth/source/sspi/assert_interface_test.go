// Copyright 2021 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package sspi_test

import (
	"creepercoding.dev/models/auth"
	"creepercoding.dev/services/auth/source/sspi"
)

// This test file exists to assert that our Source exposes the interfaces that we expect
// It tightly binds the interfaces and implementation without breaking go import cycles

type sourceInterface interface {
	auth.Config
}

var _ (sourceInterface) = &sspi.Source{}
