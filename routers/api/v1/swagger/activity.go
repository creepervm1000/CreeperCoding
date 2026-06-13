// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package swagger

import (
	api "creepercoding.dev/modules/structs"
)

// ActivityFeedsList
// swagger:response ActivityFeedsList
type swaggerActivityFeedsList struct {
	// in:body
	Body []api.Activity `json:"body"`
}
