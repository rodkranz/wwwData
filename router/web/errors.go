// Copyright 2016 Kranz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package web

import (
	"github.com/rodkranz/tmp/modules/base"
	"github.com/rodkranz/tmp/modules/context"
	"net/http"
	"fmt"
)

const (
	Error404 base.TplName = "status/NotFound"
)

func NotFound(ctx *context.Context) {
	ctx.Data["Title"] = "Page Not Found"

	ctx.Data["ErrorTitle"] = http.StatusNotFound
	ctx.Data["ErrorSmall"] = http.StatusText(http.StatusNotFound)
	ctx.Data["ErrorDescription"] = fmt.Sprintf("Page [%s] not found.", ctx.Req.URL.Path)

	ctx.HTML(404, Error404)
}

