// Copyright 2016 Kranz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package test

import (
	"net/http"

	"github.com/rodkranz/wwwData/modules/context"
)

func Hello(ctx *context.APIContext) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status_code": http.StatusText(http.StatusOK),
		"resource":    "Response json ok!",
	})
}
