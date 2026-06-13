// Copyright 2025 The CreeperCoding Authors
// Copyright 2016 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !bindata

package public

import (
	"creepercoding.dev/modules/assetfs"
	"creepercoding.dev/modules/setting"
)

func BuiltinAssets() *assetfs.Layer {
	return assetfs.Local("builtin(static)", setting.StaticRootPath, "public")
}
