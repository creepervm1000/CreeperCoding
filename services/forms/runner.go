// Copyright 2022 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package forms

import (
	"net/http"

	"creepercoding.dev/modules/web/middleware"
	"creepercoding.dev/services/context"

	"gitea.com/go-chi/binding"
)

// EditRunnerForm form for admin to create runner
type EditRunnerForm struct {
	Description string
}

// Validate validates form fields
func (f *EditRunnerForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}
