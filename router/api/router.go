// Copyright 2016 Kranz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package api

import (
	"gopkg.in/macaron.v1"
	//"github.com/go-macaron/binding"

	"github.com/rodkranz/tmp/modules/context"
	"github.com/rodkranz/tmp/router/api/v1/test"
)

// Contexter middleware already checks token for user sign in process.
func reqToken() macaron.Handler {
	return func(ctx *context.Context) {
		if !ctx.IsSigned {
			ctx.Error(401)
			return
		}
	}
}

func reqBasicAuth() macaron.Handler {
	return func(ctx *context.Context) {
		if !ctx.IsBasicAuth {
			ctx.Error(401)
			return
		}
	}
}


func RegisterRoutes(m *macaron.Macaron) {
	//bind := binding.Bind

	m.Group("/v1", func() {
		m.Get("hello", test.Hello)

	}, context.APIContexter())
}
