package main

import (
	"fmt"

	"github.com/Base-Technology/base-app-lite/common"
	"github.com/Base-Technology/base-app-lite/conf"
	_ "github.com/Base-Technology/base-app-lite/ctrl"
	"github.com/Base-Technology/base-app-lite/database"
	"github.com/Base-Technology/base-app-lite/gin"
	"github.com/Base-Technology/base-app-lite/seelog"
)

func main() {
	defer seelog.Flush()
	seelog.Infof("start base-backend-lite")

	if err := conf.InitConfig(); err != nil {
		seelog.Errorf("init config error, %v", err)
		return
	}

	if err := database.InitDatabase(); err != nil {
		seelog.Errorf("init database error, %v", err)
		return
	}

	router := gin.CreateGin()
	common.RouterRegister.SetRouter(router)
	common.RouterRegister.InitRouter()
	if err := router.Run(fmt.Sprintf(":%d", conf.Conf.ServerConf.Port)); err != nil {
		seelog.Errorf("router run error, %v", err)
		return
	}

	seelog.Infof("stop base-backend-lite")
}
