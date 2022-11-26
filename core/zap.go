package core

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"mybox/core/internal"
	"mybox/global"
	"mybox/utils"
	"os"
)

func Zap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.BOX_CONFIG.Zap.Director); !ok {
		fmt.Printf("create %v directory\n", global.BOX_CONFIG.Zap.Director)
		_ = os.Mkdir(global.BOX_CONFIG.Zap.Director, os.ModePerm)
	}

	cores := internal.Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if global.BOX_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
