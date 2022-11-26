package service

import (
	"context"
	"mybox/global"
	"mybox/utils"
)

type JwtService struct{}

func (js *JwtService) GetRedisJwt(username string) (redisJWT string, err error) {
	return global.BOX_REDIS.Get(context.Background(), username).Result()
}

func (js *JwtService) SetRedisJwt(username, token string) (err error) {
	dr, err := utils.ParseDuration(global.BOX_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	return global.BOX_REDIS.Set(context.Background(), username, token, dr).Err()

}
