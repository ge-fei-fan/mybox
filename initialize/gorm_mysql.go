package initialize

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mybox/global"
	"mybox/initialize/internal"
)

func GormMysql() *gorm.DB {
	m := global.BOX_CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}
	mysqlconfig := mysql.Config{
		DSN:                       m.Dsn(),
		DefaultStringSize:         256,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	db, err := gorm.Open(mysql.New(mysqlconfig), internal.Gorm.Config())
	if err != nil {
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
