package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mybox/config"
)

var (
	BOX_DB     *gorm.DB
	BOX_CONFIG config.Server
	BOX_VP     *viper.Viper
	BOX_LOG    *zap.Logger
	BOX_REDIS  *redis.Client
)
