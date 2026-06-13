// Copyright 2021 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package common

import (
	"path"
	"time"

	repo_model "creepercoding.dev/models/repo"
	"creepercoding.dev/modules/git"
	"creepercoding.dev/modules/httpcache"
	"creepercoding.dev/modules/httplib"
	"creepercoding.dev/modules/setting"
	"creepercoding.dev/modules/structs"
	"creepercoding.dev/services/context"
)

// ServeBlob download a git.Blob
func ServeBlob(ctx *context.Base, repo *repo_model.Repository, filePath string, blob *git.Blob, lastModified *time.Time) error {
	if httpcache.HandleGenericETagPrivateCache(ctx.Req, ctx.Resp, `"`+blob.ID.String()+`"`, lastModified) {
		return nil
	}

	if err := repo.LoadOwner(ctx); err != nil {
		return err
	}

	dataRc, err := blob.DataAsync()
	if err != nil {
		return err
	}
	defer dataRc.Close()

	if lastModified == nil {
		lastModified = new(time.Time)
	}
	httplib.ServeUserContentByReader(ctx.Req, ctx.Resp, blob.Size(), dataRc, httplib.ServeHeaderOptions{
		Filename:      path.Base(filePath),
		CacheIsPublic: !repo.IsPrivate && repo.Owner.Visibility == structs.VisibleTypePublic,
		CacheDuration: setting.StaticCacheTime,
		LastModified:  *lastModified,
	})
	return nil
}
