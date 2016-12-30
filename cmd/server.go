// Copyright 2016 Kranz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"net/http"
	"path"
	"io/ioutil"
	"fmt"
	"strings"
	"crypto/tls"
	"net/http/fcgi"
	"os"
	"net"

	"gopkg.in/macaron.v1"
	"gopkg.in/urfave/cli.v2"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/session"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/i18n"

	"github.com/rodkranz/tmp/modules/context"
	"github.com/rodkranz/tmp/modules/template"
	"github.com/rodkranz/tmp/modules/setting"
	"github.com/rodkranz/tmp/modules/log"
	"github.com/rodkranz/tmp/modules/bindata"
	"github.com/rodkranz/tmp/router"

	routerWeb "github.com/rodkranz/tmp/router/web"
	routerApi "github.com/rodkranz/tmp/router/api"
)

// checkVersion checks if binary matches the version of templates files.
func checkVersion() {
	// Templates.
	data, err := ioutil.ReadFile(setting.StaticRootPath + "/templates/.VERSION")
	if err != nil {
		log.Fatal(4, "Fail to read 'templates/.VERSION': %v", err)
	}

	if string(data) != setting.AppVer {
		log.Fatal(4, "Binary and template file version does not match, did you forget to recompile?")
	}
}

var Server = &cli.Command{
	Name:        "server",
	Usage:       "Run Server",
	Description: `Start server.`,
	Action:      runServer,
	Flags:       []cli.Flag{},
}

func newMacaron() *macaron.Macaron {
	m := macaron.New()

	// Logs
	if !setting.DisableRouterLog {
		m.Use(macaron.Logger())
	}

	// Gzip compress htmls
	m.Use(macaron.Recovery())
	if setting.EnableGzip {
		m.Use(gzip.Gziper())
	}

	// Protocol
	if setting.Protocol == setting.FCGI {
		m.SetURLPrefix(setting.AppSubUrl)
	}

	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:         path.Join("templates"),
		AppendDirectories: []string{path.Join("templates")},
		Funcs:             template.NewFuncMap(),
		IndentJSON:        macaron.Env != macaron.PROD,
	}))

	// Static folder
	m.Use(macaron.Static(
		path.Join(setting.StaticRootPath, "public"),
		macaron.StaticOptions{
			SkipLogging: setting.DisableRouterLog,
		},
	))

	localeNames, err := bindata.AssetDir("conf/locale")
	if err != nil {
		log.Fatal(4, "Fail to list locale files: %v", err)
	}
	localFiles := make(map[string][]byte)
	for _, name := range localeNames {
		localFiles[name] = bindata.MustAsset("conf/locale/" + name)
	}

	m.Use(i18n.I18n(i18n.Options{
		SubURL:          setting.AppSubUrl,
		Files:           localFiles,
		CustomDirectory: path.Join(setting.CustomPath, "conf/locale"),
		Langs:           setting.Langs,
		Names:           setting.Names,
		DefaultLang:     "en-US",
		Redirect:        true,
	}))
	m.Use(cache.Cacher(cache.Options{
		Adapter:       setting.CacheAdapter,
		AdapterConfig: setting.CacheConn,
		Interval:      setting.CacheInterval,
	}))

	m.Use(session.Sessioner(setting.SessionConfig))

	m.Use(csrf.Csrfer(csrf.Options{
		Secret:     setting.SecretKey,
		Cookie:     setting.CSRFCookieName,
		SetCookie:  true,
		Header:     "X-Csrf-Token",
		CookiePath: setting.AppSubUrl,
	}))

	m.Use(context.Contexter())
	return m
}

func runServer(ctx *cli.Context) error {
	if ctx.IsSet("config") {
		setting.CustomConf = ctx.String("config")
	}
	router.GlobalInit()

	m := newMacaron()

	// Web
	m.Get("/", routerWeb.Home)
	m.Get("/about", routerWeb.About)

	// Api
	m.Group("/api", func() {
		routerApi.RegisterRoutes(m)
	}, context.APIContexter())

	// robots.txt
	m.Get("/robots.txt", func(ctx *context.Context) {
		if setting.HasRobotsTxt {
			ctx.ServeFileContent(path.Join(setting.CustomPath, "robots.txt"))
		} else {
			ctx.Error(http.StatusNotFound)
		}
	})

	// Not found handler.
	m.NotFound(routerWeb.NotFound)

	// Flag for port number in case first time run conflict.
	if ctx.IsSet("port") {
		setting.AppUrl = strings.Replace(setting.AppUrl, setting.HTTPPort, ctx.String("port"), 1)
		setting.HTTPPort = ctx.String("port")
	}

	var listenAddr string
	if setting.Protocol == setting.UNIX_SOCKET {
		listenAddr = fmt.Sprintf("%s", setting.HTTPAddr)
	} else {
		listenAddr = fmt.Sprintf("%s:%s", setting.HTTPAddr, setting.HTTPPort)
	}
	log.Info("Listen: %v://%s%s", setting.Protocol, listenAddr, setting.AppSubUrl)

	var err error
	switch setting.Protocol {
	case setting.HTTP:
		err = http.ListenAndServe(listenAddr, m)
	case setting.HTTPS:
		server := &http.Server{Addr: listenAddr, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS10}, Handler: m}
		err = server.ListenAndServeTLS(setting.CertFile, setting.KeyFile)
	case setting.FCGI:
		err = fcgi.Serve(nil, m)
	case setting.UNIX_SOCKET:
		os.Remove(listenAddr)

		var listener *net.UnixListener
		listener, err = net.ListenUnix("unix", &net.UnixAddr{listenAddr, "unix"})
		if err != nil {
			break // Handle error after switch
		}

		// FIXME: add proper implementation of signal capture on all protocols
		// execute this on SIGTERM or SIGINT: listener.Close()
		if err = os.Chmod(listenAddr, os.FileMode(setting.UnixSocketPermission)); err != nil {
			log.Fatal(4, "Failed to set permission of unix socket: %v", err)
		}
		err = http.Serve(listener, m)
	default:
		log.Fatal(4, "Invalid protocol: %s", setting.Protocol)
	}

	if err != nil {
		log.Fatal(4, "Fail to start server: %v", err)
	}

	return nil
}
