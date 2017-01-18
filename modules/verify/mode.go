// Copyright 2017 Kranz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package verify

import (
	"strings"

	"gopkg.in/macaron.v1"

	"github.com/rodkranz/wwwData/modules/log"
	"github.com/rodkranz/wwwData/modules/setting"
)

func CheckRunMode() {
	switch setting.Cfg.Section("app").Key("RUN_MODE").String() {
	case "prod":
		macaron.Env = macaron.PROD
		macaron.ColorLog = false
		setting.ProdMode = true
	}
	log.Info("Run Mode: %s", strings.Title(macaron.Env))
}
