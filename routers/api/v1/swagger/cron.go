// Copyright 2020 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package swagger

import (
	api "creepercoding.dev/modules/structs"
)

// CronList
// swagger:response CronList
type swaggerResponseCronList struct {
	// in:body
	Body []api.Cron `json:"body"`
}
