// Copyright 2016 Kranz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package setting

import (
	"os/exec"
	"os"
	"path/filepath"
	"runtime"
	"log"
	"strings"
	"net/url"
	"path"

	"gopkg.in/ini.v1"
	"github.com/Unknwon/com"
	"github.com/go-macaron/session"
	"strconv"
)

type Scheme string

const (
	HTTP        Scheme = "http"
	HTTPS       Scheme = "https"
	FCGI        Scheme = "fcgi"
	UNIX_SOCKET Scheme = "unix"
)

var (
	// App Serrings
	AppVer         string
	AppName        string
	AppDesc        string
	AppUrl         string
	AppPath        string
	AppSubUrl      string
	AppSubUrlDepth int // Number of slashes
	AppDataPath    string

	// Server settings
	Protocol             Scheme
	Domain               string
	HTTPAddr, HTTPPort   string
	LocalURL             string
	DisableRouterLog     bool
	CertFile, KeyFile    string
	StaticRootPath       string
	EnableGzip           bool
	UnixSocketPermission uint32

	// Security settings
	SecretKey          string
	LogInRememberDays  int
	CookieUserName     string
	CookieRememberName string

	// Cache settings
	CacheAdapter  string
	CacheInterval int
	CacheConn     string

	// Session settings
	SessionConfig  session.Options
	CSRFCookieName = "_csrf"

	// Global setting objects
	Cfg          *ini.File
	CustomPath   string // Custom directory path
	CustomConf   string
	ProdMode     bool
	IsWindows    bool
	HasRobotsTxt bool

	// Log settings
	LogRootPath string
	LogModes    []string
	LogConfigs  []string

	// Api
	AllowCrossDomain bool

	// I18n settings
	Langs, Names []string
	TimeFormat   string
	dateLangs    map[string]string
)

// DateLang transforms standard language locale name to corresponding value in datetime plugin.
func DateLang(lang string) string {
	name, ok := dateLangs[lang]
	if ok {
		return name
	}
	return "en"
}

// execPath returns the executable path.
func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func init() {
	IsWindows = runtime.GOOS == "windows"
	//log.NewLogger(0, "console", `{"level": 0}`)

	var err error
	if AppPath, err = execPath(); err != nil {
		//		log.Fatal(4, "fail to get app path: %v\n", err)
	}

	// Note: we don't use path.Dir here because it does not handle case
	//      which path starts with two "/" in Windows: "//psf/Home/..."
	AppPath = strings.Replace(AppPath, "\\", "/", -1)
}

// WorkDir returns absolute path of work directory.
func WorkDir() (string, error) {
	wd := os.Getenv("WWW_DIR")
	if len(wd) > 0 {
		return wd, nil
	}

	i := strings.LastIndex(AppPath, "/")
	if i == -1 {
		return AppPath, nil
	}
	return AppPath[:i], nil
}

func forcePathSeparator(path string) {
	if strings.Contains(path, "\\") {
		log.Fatal(4, "Do not use '\\' or '\\\\' in paths, instead, please use '/' in all places")
	}
}

// NewContext initializes configuration context.
// NOTE: do not print any log except error.
func NewContext() {
	workDir, err := WorkDir()
	if err != nil {
		log.Fatal(4, "Fail to get work directory: %v", err)
	}

	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatal(4, "Fail to parse 'conf/app.ini': %v", err)
	}

	CustomPath = os.Getenv("HD_CUSTOM")
	if len(CustomPath) == 0 {
		CustomPath = workDir + "/custom"
	}

	if len(CustomConf) == 0 {
		CustomConf = CustomPath + "/conf/app.ini"
	}

	if com.IsFile(CustomConf) {
		if err = Cfg.Append(CustomConf); err != nil {
			//			log.Fatal(4, "Fail to load custom conf '%s': %v", CustomConf, err)
		}
	} else {
		//log.Warn("Custom config '%s' not found, ignore this if you're running first time", CustomConf)
	}
	Cfg.NameMapper = ini.AllCapsUnderscore

	sec := Cfg.Section("app")
	AppName = sec.Key("NAME").MustString("HD: Hey Driver Service")
	AppDesc = sec.Key("DESCRIPTION").MustString("Hey Driver App Description")

	sec = Cfg.Section("server")
	AppUrl = sec.Key("ROOT_URL").MustString("http://localhost:9090/")
	if AppUrl[len(AppUrl) - 1] != '/' {
		AppUrl += "/"
	}

	// Check if has app suburl.
	surl, err := url.Parse(AppUrl)
	if err != nil {
		//log.Fatal(4, "Invalid ROOT_URL '%s': %s", AppUrl, err)
	}

	// Suburl should start with '/' and end without '/', such as '/{subpath}'.
	AppSubUrl = strings.TrimSuffix(surl.Path, "/")
	AppSubUrlDepth = strings.Count(AppSubUrl, "/")

	Protocol = HTTP
	if sec.Key("PROTOCOL").String() == "https" {
		Protocol = HTTPS
		CertFile = sec.Key("CERT_FILE").String()
		KeyFile = sec.Key("KEY_FILE").String()
	} else if sec.Key("PROTOCOL").String() == "unix" {
		Protocol = UNIX_SOCKET
		UnixSocketPermissionRaw := sec.Key("UNIX_SOCKET_PERMISSION").MustString("666")
		UnixSocketPermissionParsed, err := strconv.ParseUint(UnixSocketPermissionRaw, 8, 32)
		if err != nil || UnixSocketPermissionParsed > 0777 {
			log.Fatal(4, "Fail to parse unixSocketPermission: %s", UnixSocketPermissionRaw)
		}
		UnixSocketPermission = uint32(UnixSocketPermissionParsed)
	}

	Domain = sec.Key("DOMAIN").MustString("localhost")
	HTTPAddr = sec.Key("HTTP_ADDR").MustString("0.0.0.0")
	HTTPPort = sec.Key("HTTP_PORT").MustString("3000")
	LocalURL = sec.Key("LOCAL_ROOT_URL").MustString("http://localhost:" + HTTPPort + "/")
	DisableRouterLog = sec.Key("DISABLE_ROUTER_LOG").MustBool()
	AppDataPath = sec.Key("APP_DATA_PATH").MustString("data")
	StaticRootPath = sec.Key("STATIC_ROOT_PATH").MustString(workDir)
	EnableGzip = sec.Key("ENABLE_GZIP").MustBool()

	HasRobotsTxt = com.IsFile(path.Join(CustomPath, "robots.txt"))
	AllowCrossDomain = Cfg.Section("api").Key("ALLOW_CROSS_DOMAIN").MustBool(true)

	// Security
	sec = Cfg.Section("security")
	SecretKey = sec.Key("SECRET_KEY").String()
	LogInRememberDays = sec.Key("LOGIN_REMEMBER_DAYS").MustInt()
	CookieUserName = sec.Key("COOKIE_USERNAME").String()
	CookieRememberName = sec.Key("COOKIE_REMEMBER_NAME").String()

	sec = Cfg.Section("i18n")
	Langs = sec.Key("LANGS").Strings(",")
	Names = sec.Key("NAMES").Strings(",")
	TimeFormat = sec.Key("DATE_FORMT").MustString("Mon Jan 2 15:04:05 MST 2006")
	dateLangs = Cfg.Section("i18n.datelang").KeysHash()
}
