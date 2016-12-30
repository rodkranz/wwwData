// Copyright 2016 Kranz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package web

import (
	"net/http"

	"github.com/rodkranz/tmp/modules/base"
	"github.com/rodkranz/tmp/modules/context"
	"github.com/rodkranz/tmp/modules/setting"
)

const (
	ABOUT_TEMPLATE base.TplName = "about"
)

func About(ctx *context.Context) {
	ctx.Data["Title"] = "About " + setting.AppName
	ctx.Data["PageIsAbout"] = true

	ctx.HTML(http.StatusOK, ABOUT_TEMPLATE)
}
