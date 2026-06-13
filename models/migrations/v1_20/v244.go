// Copyright 2023 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_20

import "creepercoding.dev/models/db"

func AddNeedApprovalToActionRun(x db.EngineMigration) error {
	/*
		New index: TriggerUserID
		New fields: NeedApproval, ApprovedBy
	*/
	type ActionRun struct {
		TriggerUserID int64 `xorm:"index"`
		NeedApproval  bool  // may need approval if it's a fork pull request
		ApprovedBy    int64 `xorm:"index"` // who approved
	}

	return x.Sync(new(ActionRun))
}
