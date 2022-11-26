package router

import (
	"github.com/gin-gonic/gin"
	v1 "mybox/api/v1"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")
	baseApi := new(v1.BaseApi)
	{
		baseRouter.POST("register", baseApi.Register)
		baseRouter.POST("login", baseApi.Login)
		//baseRouter.POST("captcha", baseApi.Captcha)
	}
}
