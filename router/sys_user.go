package router

import (
	"github.com/gin-gonic/gin"
	v1 "mybox/api/v1"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	baseApi := new(v1.BaseApi)
	{
		userRouter.GET("getUserInfo", baseApi.GetUserInfo) // 获取自身信息
	}
}
