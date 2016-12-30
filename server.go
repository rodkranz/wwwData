// Copyright 2016 Kranz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"os"
	"runtime"

	"gopkg.in/urfave/cli.v2"

	"github.com/rodkranz/tmp/cmd"
	"github.com/rodkranz/tmp/modules/setting"
)

const VER = "1.0.0"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setting.AppVer = VER
}

func main() {
	app := cli.App{
		Name: "wwwBase",
		Usage: "Base for create a webapp with go mocarron",
		Version: VER,
		Commands: []*cli.Command{
			cmd.Server,
		},
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	if len(os.Args) == 1 {
		os.Args = append(os.Args, cmd.Server.Name)
	}
	app.Run(os.Args)
}
