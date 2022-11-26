package initialize

import (
	"github.com/gin-gonic/gin"
	"mybox/global"
	"mybox/middleware"
	"mybox/router"
	"net/http"
)

func Router() *gin.Engine {
	Router := gin.Default()

	//本地存储的实际路径
	Router.StaticFS(global.BOX_CONFIG.Local.Path, http.Dir(global.BOX_CONFIG.Local.StorePath))

	//不用鉴权的路由组
	PublicGroup := Router.Group("")
	{
		//todo: 注册登录路由
		//PublicGroup.GET("/test", nil)
		router.InitBaseRouter(PublicGroup)
	}

	//需要鉴权的接口
	PrivateGroup := Router.Group("")
	PrivateGroup.Use(middleware.JwtAuth())
	{
		router.InitUserRouter(PrivateGroup)
		router.InitFileUploadAndDownloadRouter(PrivateGroup)
	}
	global.BOX_LOG.Info("router register success")
	return Router
}
