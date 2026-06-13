// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package templates

import (
	"creepercoding.dev/modules/assetfs"
	"creepercoding.dev/modules/setting"
)

func AssetFS() *assetfs.LayeredFS {
	return assetfs.Layered(CustomAssets(), BuiltinAssets())
}

func CustomAssets() *assetfs.Layer {
	return assetfs.Local("custom", setting.CustomPath, "templates")
}
