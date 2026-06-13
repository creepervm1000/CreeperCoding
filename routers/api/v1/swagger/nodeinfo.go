// Copyright 2021 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package swagger

import (
	api "creepercoding.dev/modules/structs"
)

// NodeInfo
// swagger:response NodeInfo
type swaggerResponseNodeInfo struct {
	// in:body
	Body api.NodeInfo `json:"body"`
}
