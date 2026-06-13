// Copyright 2021 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package migrations

import (
	"crypto/tls"
	"net/http"

	"creepercoding.dev/modules/hostmatcher"
	"creepercoding.dev/modules/proxy"
	"creepercoding.dev/modules/setting"
)

// NewMigrationHTTPClient returns a HTTP client for migration
func NewMigrationHTTPClient() *http.Client {
	return &http.Client{
		Transport: NewMigrationHTTPTransport(),
	}
}

// NewMigrationHTTPTransport returns a HTTP transport for migration
func NewMigrationHTTPTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: setting.Migrations.SkipTLSVerify},
		Proxy:           proxy.Proxy(),
		DialContext:     hostmatcher.NewDialContext("migration", allowList, blockList, setting.Proxy.ProxyURLFixed),
	}
}
