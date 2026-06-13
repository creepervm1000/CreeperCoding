// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The CreeperCoding Authors.
// SPDX-License-Identifier: MIT

package admin

import (
	"creepercoding.dev/models/db"
	user_model "creepercoding.dev/models/user"
	"creepercoding.dev/modules/setting"
	"creepercoding.dev/modules/structs"
	"creepercoding.dev/modules/templates"
	"creepercoding.dev/routers/web/explore"
	"creepercoding.dev/services/context"
)

const (
	tplOrgs templates.TplName = "admin/org/list"
)

// Organizations show all the organizations
func Organizations(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("admin.organizations")
	ctx.Data["PageIsAdminOrganizations"] = true

	if ctx.FormString("sort") == "" {
		ctx.SetFormString("sort", UserSearchDefaultAdminSort)
	}

	explore.RenderUserSearch(ctx, user_model.SearchUserOptions{
		Actor:           ctx.Doer,
		Types:           []user_model.UserType{user_model.UserTypeOrganization},
		IncludeReserved: true, // administrator needs to list all accounts include reserved
		ListOptions: db.ListOptions{
			PageSize: setting.UI.Admin.OrgPagingNum,
		},
		Visible: []structs.VisibleType{structs.VisibleTypePublic, structs.VisibleTypeLimited, structs.VisibleTypePrivate},
	}, tplOrgs)
}
