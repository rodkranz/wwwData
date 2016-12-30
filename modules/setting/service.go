// Copyright 2016 Kranz. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package setting

import (
	"strings"

	"github.com/rodkranz/tmp/modules/log"
	"fmt"
	"path"
	"os"
)

func NewServices() {
	newCacheService()
	newLogService()
}

func newCacheService() {
	CacheAdapter = Cfg.Section("cache").Key("ADAPTER").In("memory", []string{"memory", "redis", "memcache"})
	switch CacheAdapter {
	case "memory":
		CacheInterval = Cfg.Section("cache").Key("INTERVAL").MustInt(60)
	case "redis", "memcache":
		CacheConn = strings.Trim(Cfg.Section("cache").Key("HOST").String(), "\" ")
	default:
		log.Fatal(4, "Unknown cache adapter: %s", CacheAdapter)
	}

	log.Info("Cache Service Enabled")
}


var logLevels = map[string]string{
	"Trace":    "0",
	"Debug":    "1",
	"Info":     "2",
	"Warn":     "3",
	"Error":    "4",
	"Critical": "5",
}

func newLogService() {
	log.Info("%s %s", AppName, AppVer)

	// Get and check log mode.
	LogModes = strings.Split(Cfg.Section("log").Key("MODE").MustString("console"), ",")
	LogConfigs = make([]string, len(LogModes))
	for i, mode := range LogModes {
		mode = strings.TrimSpace(mode)
		sec, err := Cfg.GetSection("log." + mode)
		if err != nil {
			log.Fatal(4, "Unknown log mode: %s", mode)
		}

		validLevels := []string{"Trace", "Debug", "Info", "Warn", "Error", "Critical"}
		// Log level.
		levelName := Cfg.Section("log."+mode).Key("LEVEL").In(
			Cfg.Section("log").Key("LEVEL").In("Trace", validLevels),
			validLevels)
		level, ok := logLevels[levelName]
		if !ok {
			log.Fatal(4, "Unknown log level: %s", levelName)
		}

		// Generate log configuration.
		switch mode {
		case "console":
			LogConfigs[i] = fmt.Sprintf(`{"level":%s}`, level)
		case "file":
			logPath := sec.Key("FILE_NAME").MustString(path.Join(LogRootPath, "gogs.log"))
			if err = os.MkdirAll(path.Dir(logPath), os.ModePerm); err != nil {
				panic(err.Error())
			}

			LogConfigs[i] = fmt.Sprintf(
				`{"level":%s,"filename":"%s","rotate":%v,"maxlines":%d,"maxsize":%d,"daily":%v,"maxdays":%d}`, level,
				logPath,
				sec.Key("LOG_ROTATE").MustBool(true),
				sec.Key("MAX_LINES").MustInt(1000000),
				1<<uint(sec.Key("MAX_SIZE_SHIFT").MustInt(28)),
				sec.Key("DAILY_ROTATE").MustBool(true),
				sec.Key("MAX_DAYS").MustInt(7))
		case "conn":
			LogConfigs[i] = fmt.Sprintf(`{"level":%s,"reconnectOnMsg":%v,"reconnect":%v,"net":"%s","addr":"%s"}`, level,
				sec.Key("RECONNECT_ON_MSG").MustBool(),
				sec.Key("RECONNECT").MustBool(),
				sec.Key("PROTOCOL").In("tcp", []string{"tcp", "unix", "udp"}),
				sec.Key("ADDR").MustString(":7020"))
		case "database":
			LogConfigs[i] = fmt.Sprintf(`{"level":%s,"driver":"%s","conn":"%s"}`, level,
				sec.Key("DRIVER").String(),
				sec.Key("CONN").String())
		}

		log.NewLogger(Cfg.Section("log").Key("BUFFER_LEN").MustInt64(10000), mode, LogConfigs[i])
		log.Info("Log Mode: %s(%s)", strings.Title(mode), levelName)
	}
}
