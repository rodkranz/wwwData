// Copyright 2017 Kranz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package verify

import (
	"io/ioutil"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/i18n"
	"github.com/go-macaron/session"
	"github.com/go-xorm/xorm"
	"github.com/mcuadros/go-version"
	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"

	"github.com/rodkranz/wwwData/modules/log"
	"github.com/rodkranz/wwwData/modules/setting"
)

type VerChecker struct {
	ImportPath string
	Version    func() string
	Expected   string
}

// checkVersion checks if binary matches the version of templates files.
func CheckVersion() {
	// checkVersion checks if binary matches the version of templates files.
	data, err := ioutil.ReadFile(setting.StaticRootPath + "/templates/.VERSION")
	if err != nil {
		log.Fatal(4, "Fail to read 'templates/.VERSION': %v", err)
	}
	tplVer := string(data)
	if tplVer != setting.AppVer {
		if version.Compare(tplVer, setting.AppVer, ">") {
			log.Fatal(4, "Binary version is lower than template file version, did you forget to recompile Gogs?")
		} else {
			log.Fatal(4, "Binary version is higher than template file version, did you forget to update template files?")
		}
	}

	/**



	 */
	// Check dependency version.
	checkers := []VerChecker{
		{"github.com/go-xorm/xorm", func() string { return xorm.Version }, "0.5.5"},
		{"github.com/go-macaron/cache", cache.Version, "0.1.2"},
		{"github.com/go-macaron/csrf", csrf.Version, "0.1.0"},
		{"github.com/go-macaron/i18n", i18n.Version, "0.3.0"},
		{"github.com/go-macaron/session", session.Version, "0.1.6"},
		{"gopkg.in/ini.v1", ini.Version, "1.8.4"},
		{"gopkg.in/macaron.v1", macaron.Version, "1.1.7"},
	}

	for _, c := range checkers {
		if !version.Compare(c.Version(), c.Expected, ">=") {
			log.Fatal(4, `Dependency outdated!
Package '%s' current version (%s) is below requirement (%s),
please use following command to update this package and recompile:
go get -u %[1]s`, c.ImportPath, c.Version(), c.Expected)
		}
	}
}
