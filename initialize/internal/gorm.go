package internal

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"mybox/global"
	"os"
	"time"
)

type DBBASE interface {
	GetLogMode() string
}

type _gorm struct {
}

var Gorm = new(_gorm)

func (g *_gorm) Config() *gorm.Config {
	config := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	log.New(os.Stdout, "\r\n", log.LstdFlags)
	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})

	var logMode DBBASE
	//先只支持mysql
	logMode = &global.BOX_CONFIG.Mysql

	switch logMode.GetLogMode() {
	case "silent", "Silent":
		config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = _default.LogMode(logger.Info)
	default:
		config.Logger = _default.LogMode(logger.Info)
	}
	return config
}
