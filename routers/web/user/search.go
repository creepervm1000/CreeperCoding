// Copyright 2022 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"net/http"

	"creepercoding.dev/models/db"
	user_model "creepercoding.dev/models/user"
	"creepercoding.dev/modules/optional"
	"creepercoding.dev/modules/setting"
	"creepercoding.dev/services/context"
	"creepercoding.dev/services/convert"
)

// SearchCandidates searches candidate users for dropdown list
func SearchCandidates(ctx *context.Context) {
	searchUserTypes := []user_model.UserType{user_model.UserTypeIndividual}
	if ctx.FormBool("orgs") {
		searchUserTypes = append(searchUserTypes, user_model.UserTypeOrganization)
	}
	users, _, err := user_model.SearchUsers(ctx, user_model.SearchUserOptions{
		Actor:       ctx.Doer,
		Keyword:     ctx.FormTrim("q"),
		Types:       searchUserTypes,
		IsActive:    optional.Some(true),
		ListOptions: db.ListOptions{PageSize: setting.UI.MembersPagingNum},
	})
	if err != nil {
		ctx.ServerError("Unable to search users", err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"data": convert.ToUsers(ctx, ctx.Doer, users)})
}
