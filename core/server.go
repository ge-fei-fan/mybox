package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mybox/global"
	"mybox/initialize"
	"net/http"
	"time"
)

func initServer(address string, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func RunServer() {
	if global.BOX_CONFIG.System.UseRedis {
		initialize.Redis()
	}

	Router := initialize.Router()

	address := fmt.Sprintf(":%d", global.BOX_CONFIG.System.Addr)
	s := initServer(address, Router)

	time.Sleep(10 * time.Microsecond)
	global.BOX_LOG.Info("server run success on ", zap.String("address", address))

	err := s.ListenAndServe()
	if err != nil {
		global.BOX_LOG.Error(err.Error())
	}

}
