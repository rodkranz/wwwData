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
	HOME_TEMPLATE base.TplName = "home"
)

func Home(ctx *context.Context) {
	ctx.Data["Title"] = "Home " + setting.AppName
	ctx.Data["PageIsHome"] = true

	ctx.HTML(http.StatusOK, HOME_TEMPLATE)
}
