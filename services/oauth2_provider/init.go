// Copyright 2024 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package oauth2_provider

import (
	"context"

	"creepercoding.dev/modules/setting"
)

// Init initializes the oauth source
func Init(ctx context.Context) error {
	if !setting.OAuth2.Enabled {
		return nil
	}

	return InitSigningKey()
}
