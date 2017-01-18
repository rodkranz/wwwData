// Copyright 2016 Kranz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package template

import (
	"html/template"
	"runtime"
	"strings"
	"time"

	"fmt"
	"github.com/rodkranz/wwwData/modules/setting"
)

func NewFuncMap() []template.FuncMap {
	return []template.FuncMap{map[string]interface{}{
		"GoVer": func() string {
			return strings.Title(runtime.Version())
		},
		"UseHTTPS": func() bool {
			return strings.HasPrefix(setting.AppUrl, "https")
		},
		"AppName": func() string {
			return setting.AppName
		},
		"AppDesc": func() string {
			return setting.AppDesc
		},
		"AppSubUrl": func() string {
			return setting.AppSubUrl
		},
		"AppUrl": func() string {
			return setting.AppUrl
		},
		"AppVer": func() string {
			return setting.AppVer
		},
		"AppDomain": func() string {
			return setting.Domain
		},
		"LoadTimes": func(startTime time.Time) string {
			return fmt.Sprint(time.Since(startTime).Nanoseconds()/1e6) + "ms"
		},
		"ToLower":  strings.ToLower,
		"DateTime": DateTime,
	}}
}

func DateTime(date time.Time) string {
	return date.Format(setting.TimeFormat)
}

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}

func Range(l int) []int {
	return make([]int, l)
}
