// Copyright 2025 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import "creepercoding.dev/modules/setting/config"

type CcopilotStruct struct {
	Enabled   *config.Option[bool]
	APIKey    *config.Option[string]
	Endpoint  *config.Option[string]
	ModelName *config.Option[string]
	MaxTokens *config.Option[int64]
}
