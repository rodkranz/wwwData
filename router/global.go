package routers

import (
	"github.com/rodkranz/wwwData/models"
	"github.com/rodkranz/wwwData/modules/log"
	"github.com/rodkranz/wwwData/modules/setting"
	"github.com/rodkranz/wwwData/modules/verify"
)

func GlobalInit() {
	setting.NewContext()

	log.Trace("Custom path: %s", setting.CustomPath)
	log.Trace("Log path: %s", setting.LogRootPath)

	models.LoadConfigs()
	setting.NewServices()

	if err := models.NewEngine(); err != nil {
		log.Fatal(4, "Fail to initialize ORM engine: %v", err)
	}
	models.HasEngine = true

	verify.CheckRunMode()
}
