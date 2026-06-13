// Copyright 2025 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package label

import (
	"creepercoding.dev/modules/label"
	"creepercoding.dev/modules/web"
	"creepercoding.dev/services/context"
	"creepercoding.dev/services/forms"
)

func GetLabelEditForm(ctx *context.Context) *forms.CreateLabelForm {
	form := web.GetForm(ctx).(*forms.CreateLabelForm)
	if ctx.HasError() {
		ctx.JSONError(ctx.GetErrMsg())
		return nil
	}
	var err error
	form.Color, err = label.NormalizeColor(form.Color)
	if err != nil {
		ctx.JSONError(ctx.Tr("repo.issues.label_color_invalid"))
		return nil
	}
	return form
}
