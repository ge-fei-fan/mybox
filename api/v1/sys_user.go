package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"mybox/global"
	"mybox/model/common/response"
	"mybox/model/system"
	systemReq "mybox/model/system/request"
	systemResp "mybox/model/system/response"
	"mybox/utils"
)

type BaseApi struct{}

func (b *BaseApi) Register(c *gin.Context) {
	var r systemReq.Register
	err := c.ShouldBind(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.RegisterVerify(r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := system.SysUser{Username: r.Username, Password: r.Password, Email: r.Username, HeaderImg: r.HeaderImg}
	userReturn, err := userService.Register(&user)
	if err != nil {
		global.BOX_LOG.Error("注册失败!", zap.Error(err))
		response.FailWithDetailed(systemResp.SysUserResponse{User: *userReturn}, "注册失败", c)
		return
	}
	response.OkWithDetailed(systemResp.SysUserResponse{User: *userReturn}, "注册成功", c)
}

func (b *BaseApi) Login(c *gin.Context) {
	var l systemReq.Login

	err := c.ShouldBind(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//todo 验证码校验
	//账号密码非空
	if len(l.Username) == 0 || len(l.Password) == 0 {
		response.FailWithMessage("用户名或者密码的值为空", c)
		return
	}
	su := system.SysUser{Username: l.Username, Password: l.Password}
	//校验用户名密码
	user, err := userService.Login(&su)
	if err != nil {
		global.BOX_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
		response.FailWithMessage("登陆失败! 用户名不存在或者密码错误!", c)
		return
	}
	//分发token
	b.GeneralToken(c, user)
}

func (b *BaseApi) GeneralToken(c *gin.Context, user *system.SysUser) {
	j := utils.NewJwt()
	claims := j.CreateClaims(&systemReq.BaseClaims{
		UUID:     user.UUID,
		ID:       user.ID,
		Username: user.Username,
	})

	token, err := j.CreateToken(claims)
	if err != nil {
		global.BOX_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("生成token失败", c)
		return
	}

	_, err = jwtService.GetRedisJwt(user.Username)
	//redis没有存储token
	if err == redis.Nil {
		err := jwtService.SetRedisJwt(user.Username, token)
		if err != nil {
			global.BOX_LOG.Error("设置登录状态失败!", zap.Error(err))
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(systemResp.LoginResponse{
			User:      *user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt,
		}, "登录成功", c)
	} else if err != nil { //设置redis报错
		global.BOX_LOG.Error("设置登录状态失败!", zap.Error(err))
		response.FailWithMessage("设置登录状态失败", c)
		return
	} else {
		//redis存在token  应该是重新设置一下
		//todo 可能需要把旧的token给拉黑取消掉
		if err := jwtService.SetRedisJwt(user.Username, token); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(systemResp.LoginResponse{
			User:      *user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt,
		}, "登录成功", c)
	}
}

func (b *BaseApi) GetUserInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	ReqUser, err := userService.GetUserInfo(uuid)
	if err != nil {
		global.BOX_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "获取成功", c)
}
