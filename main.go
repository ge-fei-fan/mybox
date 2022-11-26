package main

import (
	"go.uber.org/zap"
	"mybox/core"
	"mybox/global"
	"mybox/initialize"
)

func main() {
	global.BOX_VP = core.Viper()
	//fmt.Println(global.BOX_CONFIG.Mysql)
	global.BOX_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.BOX_LOG)
	global.BOX_DB = initialize.Gorm() //初始化数据库
	if global.BOX_DB != nil {
		initialize.RegisterTables(global.BOX_DB)
		db, _ := global.BOX_DB.DB()
		defer db.Close()
	}
	//启动gin服务
	core.RunServer()
}
