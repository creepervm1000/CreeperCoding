// Copyright 2026 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"creepercoding.dev/modules/log"
	"creepercoding.dev/modules/migration"
	"creepercoding.dev/modules/structs"
)

func init() {
	RegisterDownloaderFactory(&ForgejoDownloaderFactory{})
}

// ForgejoDownloaderFactory defines a forgejo downloader factory
type ForgejoDownloaderFactory struct{}

// New returns a Downloader related to this factory according MigrateOptions
func (f *ForgejoDownloaderFactory) New(ctx context.Context, opts migration.MigrateOptions) (migration.Downloader, error) {
	u, err := url.Parse(opts.CloneAddr)
	if err != nil {
		return nil, err
	}

	baseURL := u.Scheme + "://" + u.Host
	repoNameSpace := strings.TrimPrefix(u.Path, "/")
	repoNameSpace = strings.TrimSuffix(repoNameSpace, ".git")

	path := strings.Split(repoNameSpace, "/")
	if len(path) < 2 {
		return nil, fmt.Errorf("invalid path: %s", repoNameSpace)
	}

	repoPath := strings.Join(path[len(path)-2:], "/")
	if len(path) > 2 {
		subPath := strings.Join(path[:len(path)-2], "/")
		baseURL += "/" + subPath
	}

	log.Trace("Create forgejo downloader. BaseURL: %s RepoName: %s", baseURL, repoNameSpace)

	return NewGiteaDownloader(ctx, baseURL, repoPath, opts.AuthUsername, opts.AuthPassword, opts.AuthToken)
}

// GitServiceType returns the type of git service
func (f *ForgejoDownloaderFactory) GitServiceType() structs.GitServiceType {
	return structs.ForgejoService
}
