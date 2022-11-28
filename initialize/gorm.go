package initialize

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mybox/global"
	"mybox/model/system"
	"os"
)

func Gorm() *gorm.DB {
	return GormMysql()
}

func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		system.SysUser{},
		system.FileUploadAndDownload{},
		system.SysUserRepository{},
		system.SharePool{},
	)
	if err != nil {
		global.BOX_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.BOX_LOG.Info("register table success")
}
